package main

import (
    "fmt"
    "io"
    "os"
    "net/http"
    "path/filepath"
    "os/exec"
    "math/rand"
    "time"
    "errors"
)

func healthRequest(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

type 

// Handle a file upload to the /upload endpoint. The upload should be a multipart form with 
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    //parse the multipart form in the request
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }

    //get the file
    file, _, err := r.FormFile("file")
    if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }

    //Strip extension if there is one
    name := r.FormValue("name")
    extension := filepath.Ext(name)
    name = name[0:len(name)-len(extension)]

    // Add .wav extension
    name += ".wav"

    // Write to temp file
    tempFilePath := "/tmp/upload/" + randomString(16)
    tmpFile, err := os.Create(tempFilePath)
    defer tmpFile.Close()
    defer os.Remove(tempFilePath)
    if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }

    _, err = io.Copy(tmpFile, file);
    if  err != nil {
   	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
    }
    
    // Transccode the uploaded file
    err = resample(tempFilePath, "/data/" + name)
    if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
    }
    fmt.Fprintf(w, "success")
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w,r, "/data/" + r.FormValue("name"))
}

func init() {
    os.Mkdir("/data/", 777)
    os.Mkdir("/tmp/upload/", 777)
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/download", downloadHandler)
    http.HandleFunc("/live", healthRequest)
    http.ListenAndServe(":80", nil)
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
