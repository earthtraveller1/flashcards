package main

import (
    "net/http"
    "encoding/json"
    "io"
    "fmt"
    "os"
)

func apiCardStacks(pWriter http.ResponseWriter, pRequest *http.Request) {
	if pRequest.Method == "GET" {
		pWriter.Header().Set("Content-Type", "application/json")

		jsonCards, error := json.Marshal(globalCards)
		if error != nil {
			pWriter.WriteHeader(500)
			return
		}

		pWriter.Write(jsonCards)
	} else if pRequest.Method == "POST" {
        requestBodyRaw := make([]byte, 0, 256)

        reading := true
        for reading {
            tempBuffer := make([]byte, 256)
            bytesRead, error := pRequest.Body.Read(tempBuffer)
            if error == io.EOF {
                reading = false
            } else if error != nil {
                fmt.Fprintf(os.Stderr, "[ERROR]: %s", error)
                return
            }

            requestBodyRaw = append(requestBodyRaw, tempBuffer[:bytesRead]...)
        }

        var requestInfo CardStack
        error := json.Unmarshal(requestBodyRaw, &requestInfo)
        if error != nil {
            pWriter.WriteHeader(400)
            pWriter.Header().Set("Content-Type", "application/json")
            pWriter.Write([]byte(fmt.Sprintf("\"error\": \"%s\"", error)))

            return
        }

        globalCards = append(globalCards, requestInfo)
	}
}
