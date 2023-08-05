package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func apiCardStacks(pWriter http.ResponseWriter, pRequest *http.Request) {
	if pRequest.Method == "GET" {
		pWriter.Header().Set("Content-Type", "application/json")

		jsonCards, error := json.Marshal(globalCardStacks)
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

		stackId := strings.ReplaceAll(requestInfo.Name, " ", "_")
		stackId = strings.ToLower(stackId)

		globalCardStacks[stackId] = requestInfo
	} else {
        pWriter.WriteHeader(405)
        pWriter.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(pWriter, `{ "error": "method not allowed" }`)
    }
}

func apiSpecificCardStack(pWriter http.ResponseWriter, pRequest *http.Request) {
	uriParts := strings.Split(pRequest.URL.RequestURI(), "/")
	stackID := uriParts[len(uriParts)-1]

	if pRequest.Method == "GET" {
		pWriter.Header().Set("Content-Type", "application/json")
		requestedStack, exists := globalCardStacks[stackID]

		if !exists {
			pWriter.WriteHeader(404)
			pWriter.Write([]byte(fmt.Sprintf(`{ "error": "the %s stack does not exist"}`, stackID)))

			return
		}

		stackJSON, error := json.Marshal(requestedStack)

		if error != nil {
			pWriter.WriteHeader(500)
			pWriter.Write([]byte(`{ "error": "internal server error"}`))

			return
		}

		pWriter.Write(stackJSON)
	} else if pRequest.Method == "DELETE" {
        _, stackExists := globalCardStacks[stackID]
        
        if !stackExists {
            pWriter.Header().Set("Content-Type", "application/json")
            pWriter.WriteHeader(404)
            fmt.Fprintf(pWriter, `{ "error": "the %s stack does not exist" }`, stackID)

            return
        }

        delete(globalCardStacks, stackID)
    } else {
        pWriter.WriteHeader(405)
        pWriter.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(pWriter, `{ "error": "method not allowed" }`)
    }

}
