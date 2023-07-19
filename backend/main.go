package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type CardStack struct {
    Name string
    Description string
}

var globalCards []CardStack

func staticFiles(pWriter http.ResponseWriter, pRequest *http.Request) {
    fileName := pRequest.URL.RequestURI()
    if fileName == "/" {
        fileName = "index.html"
    }

    if strings.HasSuffix(fileName, ".css") {
        pWriter.Header().Set("Content-Type", "text/css")
    }

    fileContent, error := ioutil.ReadFile("../frontend/" + fileName)
    if error != nil {
        pWriter.Write([]byte("<h1>404</h1> <p>Not Found"))
    }

    pWriter.Write(fileContent)
}

func doesFileExist(pFilename string) bool {
    _, error := os.Stat(pFilename)
    return !os.IsNotExist(error)
}

func getCardStacks(pWriter http.ResponseWriter, pRequest *http.Request) {
    if pRequest.Method == "GET" {
        pWriter.Header().Set("Content-Type", "application/json")

        jsonCards, error := json.Marshal(globalCards)
        if error != nil {
            pWriter.WriteHeader(500)
            return
        } 

        pWriter.Write(jsonCards)
    }
}

func main() {
    globalCards = append(globalCards, CardStack { Name: "Neng", Description: "Neng Li is the President of China" })
    globalCards = append(globalCards, CardStack { Name: "Prussia", Description: "German state during the 1800s or something." })

    http.HandleFunc("/", staticFiles)
    http.HandleFunc("/api/cardstacks", getCardStacks)

    error := http.ListenAndServe("127.0.0.1:3000", nil)
    if error != nil {
        fmt.Fprintf(os.Stderr, "[ERROR]: The server encountered an error: %s", error)
    }
}
