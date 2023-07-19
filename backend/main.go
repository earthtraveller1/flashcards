package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func staticFiles(pWriter http.ResponseWriter, pRequest *http.Request) {
    fileContent, error := ioutil.ReadFile("../frontend/" + pRequest.URL.RequestURI())
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
