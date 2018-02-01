package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "time"
    "errors"
    "net/http"
    "github.com/Shopify/sarama"
    "github.com/movio/kasper"
    "strconv"
)

var outputTopic string

func main() {
    maxMessageSize := os.Getenv("MAX_MESSAGE_SIZE")
    kafkaAddress := os.Getenv("KAFKA_BROKER_URL")
    loggingTopic := os.Getenv("LOGGING_TOPIC")
    inputTopic := os.Getenv("INPUT_TOPIC")
    outputTopic = os.Getenv("OUTPUT_TOPIC")

    _, atoi_err := strconv.Atoi(maxMessageSize)
    if maxMessageSize == "" || kafkaAddress == "" || loggingTopic == "" ||
    	inputTopic == "" || outputTopic == "" || atoi_err != nil {
	    panic("Config not correct: ")
	}

    maxMessageSizeInt, atoi_err := strconv.ParseInt(maxMessageSize, 10, 16)
    sarama.MaxRequestSize = int32(maxMessageSizeInt)
    sarama.MaxResponseSize = sarama.MaxRequestSize
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.MaxMessageBytes = int(sarama.MaxRequestSize)
    saramaClient, err := sarama.NewClient([]string{kafkaAddress}, saramaConfig)
    if err != nil {
	panic(err)
    }

    kasperConfig := &kasper.Config{
	// Topic used for logging
	TopicProcessorName:    "logging.transcoder",
	Client:                saramaClient,
	InputTopics:           []string{inputTopic},
	InputPartitions:       []int{0},
	BatchSize:             10,
	BatchWaitDuration:     1 * time.Second,
	Logger:                kasper.NewJSONLogger("transcoder", true),
	MetricsProvider:       kasper.NewPrometheus("transcoder"),
	MetricsUpdateInterval: 60 * time.Second,
    }

    messageProcessorMap := map[int]kasper.MessageProcessor{
	0 : &AudioTranscoder{}}

    topicProcessor := kasper.NewTopicProcessor(kasperConfig, messageProcessorMap)
    err = topicProcessor.RunLoop()
    log.Printf("Topic processor finished with err = %s\n", err)
}

func healthRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

func resample(input string, output string) error {
    binary, err := exec.LookPath("ffmpeg")
    if err != nil {
	panic(err)
    }

    args := []string{
	"-i", input,
	"-loglevel", "error",
	"-y",
	"-ac", "1", "-ar", "8000",
	"-acodec", "pcm_s16le",
	"-f", "wav",
        output }

    cmd := exec.Command(binary, args...)
    cmd.Env = os.Environ()

    combinedOutput, err := cmd.CombinedOutput()
    if err != nil {
	return errors.New(string(combinedOutput))
    }
    return nil
}


type AudioTranscoder struct {
}

func (*AudioTranscoder) Process(messages []*sarama.ConsumerMessage, sender kasper.Sender) error {
    log.Println("Transcoding")
    dir, err := ioutil.TempDir("","transcoding")
    if err != nil {
	panic(err)
    }
    defer os.RemoveAll(dir)

    for _, message := range messages {
	log.Println("Recieved message")
	tmpInputFile, err := ioutil.TempFile(dir, "input")
	if err != nil {
	    panic(err)
	}
	defer os.Remove(tmpInputFile.Name())
	err = ioutil.WriteFile(tmpInputFile.Name(), message.Value, os.FileMode(0777))
	if err != nil {
	    panic(err)
	}

	tmpOutputFile, err := ioutil.TempFile(dir, "output")
	if err != nil {
	    panic(err)
	}
	defer os.Remove(tmpOutputFile.Name())

	err = resample(tmpInputFile.Name(), tmpOutputFile.Name())
	if err != nil {
		// Do something idk
		panic(err)
	}

	buf, err := ioutil.ReadFile(tmpOutputFile.Name())
	if err != nil {
	    panic(err)
	}
	log.Printf("Going to send message, size: %v",len(buf))
	out := sarama.ProducerMessage{
	    Key : sarama.ByteEncoder(message.Key),
	    Value : sarama.ByteEncoder(buf),
	    Topic : outputTopic}
	sender.Send(&out)
    }
    errs := sender.Flush()
    if errs != nil {
	for _, err := range errs.(sarama.ProducerErrors) {
		log.Println("Write to kafka failed: ", err)
	}
    }
    return errs


}
