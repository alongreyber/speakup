package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

type UploadedFile struct {
    User string
    Id string
    Extension string
    Data []byte
}

func (f *UploadedFile) save() error {
    filename := "/data/" + f.User + "/" + f.Id + "." + f.Extension
    return ioutil.writeFile(filename, f.Data, 0600)
}

func read(user string, id string, extension string) *UploadedFile {
    filename := "/data/" + user + "/" + id + "." + extension
    data, _ := ioutil.ReadFile(filename)
    return &UploadedFile{User: user, Id: id, Extension: extension, Data: data}
}

func handler(w http.ResponseWriter, r *http.Request) {

}

func main() {
    fmt.Printf("Hello, world.\n")

}
