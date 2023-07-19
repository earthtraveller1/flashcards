package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func staticFiles(pWriter http.ResponseWriter, pRequest *http.Request) {
    fileName := pRequest.URL.RequestURI()
    if fileName == "/" {
        fileName = "index.html"
    }

    fileContent, error := ioutil.ReadFile("../frontend/" + fileName)
    if error != nil {
        pWriter.Write([]byte("<h1>404</h1> <p>Not Found"))
    }

    pWriter.Write(fileContent)
}

func main() {
    http.HandleFunc("/", staticFiles)

    error := http.ListenAndServe("127.0.0.1:3000", nil)
    if error != nil {
        fmt.Fprintf(os.Stderr, "[ERROR]: The server encountered an error: %s", error)
    }
}
