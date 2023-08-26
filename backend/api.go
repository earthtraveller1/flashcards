package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type StackNotExistError struct{}

func (error StackNotExistError) Error() string {
	return "the requested stack does not exists"
}

func getCardStacks(pWriter http.ResponseWriter, pRequest *http.Request, pStacks *map[string]CardStack) {
	pWriter.Header().Set("Content-Type", "application/json")

	jsonCards, error := json.Marshal(*pStacks)
	if error != nil {
		pWriter.WriteHeader(500)
		return
	}

	pWriter.Write(jsonCards)
}

func createCardStack(pJSON []byte, pWriter http.ResponseWriter, pStacks *map[string]CardStack) {
	var requestInfo CardStack
	error := json.Unmarshal(pJSON, &requestInfo)
	if error != nil {
		pWriter.WriteHeader(400)
		pWriter.Header().Set("Content-Type", "application/json")
		pWriter.Write([]byte(fmt.Sprintf("\"error\": \"%s\"", error)))

		return
	}

	stackId := strings.ReplaceAll(requestInfo.Name, " ", "_")
	stackId = strings.ToLower(stackId)

	(*pStacks)[stackId] = requestInfo
}

func addCardToStack(pCard Card, pStackName string, pStacks *map[string]CardStack) error {
	stack, stackExists := (*pStacks)[pStackName]
	if stackExists {
		stack.Cards = append(stack.Cards, pCard)
		(*pStacks)[pStackName] = stack

		return nil
	}

	return StackNotExistError{}
}

func removeFromStack(pStack *CardStack, pIndex int) {
	if pIndex > len(pStack.Cards) {
		return
	}

	(*pStack).Cards[pIndex] = (*pStack).Cards[len((*pStack).Cards)-1]
	(*pStack).Cards = (*pStack).Cards[:len((*pStack).Cards)-1]
}

func apiCardStacksHandler(pWriter http.ResponseWriter, pRequest *http.Request) {
	if pRequest.Method == "GET" {
		getCardStacks(pWriter, pRequest, &globalCardStacks)
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

		createCardStack(requestBodyRaw, pWriter, &globalCardStacks)
	} else {
		pWriter.WriteHeader(405)
		pWriter.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(pWriter, `{ "error": "method not allowed" }`)
	}
}

func apiCardStacksCardsHandler(pStackID string, uriParts []string, pWriter http.ResponseWriter, pRequest *http.Request) {
	if pRequest.Method == "POST" {
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

		var card Card
		error := json.Unmarshal(requestBodyRaw, &card)
		if error != nil {
			fmt.Fprintf(pWriter, `{ "error": "%s" }`, error)
			pWriter.Header().Set("Content-Type", "application/json")
			pWriter.WriteHeader(400)

			return
		}

		addCardToStack(card, pStackID, &globalCardStacks)
	} else if pRequest.Method == "DELETE" {
		if len(uriParts) < 6 {
			fmt.Fprintf(pWriter, `{ "error": "method not allowed" }`)
			pWriter.Header().Set("Content-Type", "application/json")
			pWriter.WriteHeader(405)

			return
		}

		cardIndex, err := strconv.Atoi(uriParts[5])
		if err != nil {
			fmt.Fprintf(pWriter, `{ "error": "invalid request" }`)
			pWriter.Header().Set("Content-Type", "application/json")
			pWriter.WriteHeader(400)

			return
		}

		stack := globalCardStacks[pStackID]

		removeFromStack(&stack, cardIndex)

		globalCardStacks[pStackID] = stack
	} else {
		fmt.Fprintf(pWriter, `{ "error": "method not allowed" }`)
		pWriter.Header().Set("Content-Type", "application/json")
		pWriter.WriteHeader(405)
	}
}

func apiSpecificCardStackHandler(pWriter http.ResponseWriter, pRequest *http.Request) {
	uriParts := strings.Split(pRequest.URL.RequestURI(), "/")
	stackID := uriParts[3]

	// Reroute to the cards handler if there is more to the URI
	if len(uriParts) >= 5 && uriParts[4] == "cards" {
		apiCardStacksCardsHandler(stackID, uriParts, pWriter, pRequest)
		return
	}

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
