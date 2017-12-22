package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "math/rand"
    "time"
    "errors"
    "net/http"
    "github.com/Shopify/sarama"
    "github.com/movio/kasper"
)

func main() {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
    saramaClient, err := sarama.NewClient([]string{"kafka-0.kafka-svc"}, saramaConfig)
    if err != nil {
	panic(err)
    }

    kasperConfig := &kasper.Config{
	// Topic used for logging
	TopicProcessorName:    "logging.transcoder",
	Client:                saramaClient,
	InputTopics:           []string{"streaming.transcoder.uploaded"},
	InputPartitions:       []int{0, 1},
	BatchSize:             10,
	BatchWaitDuration:     1 * time.Second,
	Logger:                kasper.NewJSONLogger("transcoder", false),
	MetricsProvider:       kasper.NewPrometheus("transcoder"),
	MetricsUpdateInterval: 60 * time.Second,
    }

    rand.Seed(time.Now().UTC().UnixNano())

    // One for each partition
    messageProcessorMap := make(map[int]kasper.MessageProcessor)
    i := 0
    for _, _ = range messageProcessorMap {
	messageProcessorMap[i] = &AudioTranscoder{}
	i++
    }

    topicProcessor := kasper.NewTopicProcessor(kasperConfig, messageProcessorMap)
    topicProcessor.RunLoop()
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
    dir, err := ioutil.TempDir("","transcoding")
    if err != nil {
	panic(err)
    }
    defer os.RemoveAll(dir)

    for _, message := range messages {
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
	out := sarama.ProducerMessage{
	    Key : sarama.ByteEncoder(message.Key),
	    Value : sarama.ByteEncoder(buf),
	    Topic : "streaming.transcoder.transcoded"}
	sender.Send(&out)
    }
    return nil

}
func randomString(l int) string {
    bytes := make([]byte, l)
    for i := 0; i < l; i++ {
        bytes[i] = byte(randInt(97, 122))
    }
    return string(bytes)
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}
