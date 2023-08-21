package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type CardSlice []Card

type CardStack struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cards       CardSlice `json:"cards"`
}

func (slice CardSlice) MarshalJSON() ([]byte, error) {
	if len(slice) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal([]Card(slice))
}

var globalCardStacks map[string]CardStack

func staticFilesHandler(pWriter http.ResponseWriter, pRequest *http.Request) {
	fileName := pRequest.URL.RequestURI()
	if fileName == "/" {
		fileName = "index.html"
	}

	if strings.HasSuffix(fileName, ".css") {
		pWriter.Header().Set("Content-Type", "text/css")
	}

	if strings.HasSuffix(fileName, ".js") {
		pWriter.Header().Set("Content-Type", "application/javascript")
	}

	fileContent, error := os.ReadFile("../frontend/" + fileName)
	if error != nil {
		pWriter.Write([]byte("<h1>404</h1> <p>Not Found"))
	}

	pWriter.Write(fileContent)
}

func doesFileExist(pFilename string) bool {
	_, error := os.Stat(pFilename)
	return !os.IsNotExist(error)
}

func stackPage(pTemplate, pStackName string) string {
	subsitution := "//{{{server}}}"
	replacement := fmt.Sprintf(`
const serverInfo = {
    stackName: "%s"
}
    `, pStackName)

	if strings.Contains(pTemplate, subsitution) {
		return strings.Replace(pTemplate, subsitution, replacement, 1)
	}

	return pTemplate
}

func stackPageHander(pWriter http.ResponseWriter, pRequest *http.Request) {
	uri := pRequest.URL.RequestURI()
	uriParts := strings.Split(uri, "/")

	// The first element is blank because of how the string split function works.
	stackName := uriParts[2]

	_, stackExists := globalCardStacks[stackName]
	if !stackExists {
		pWriter.WriteHeader(404)
		fmt.Fprintf(pWriter, "<h1>404</h1> <p>Not found")
		return
	}

	stackPageTemplate, error := os.ReadFile("stack.html")
	if error != nil {
		pWriter.WriteHeader(500)
		pWriter.Write([]byte("<h1>500</h1> <p>Internal server error"))

		fmt.Fprintf(os.Stderr, "[ERROR]: stack.html does not appear to exist, or I cannot load it or something.\n")
	}

	pWriter.Write([]byte(stackPage(string(stackPageTemplate), stackName)))
}

func runServer(waitGroup *sync.WaitGroup, server *http.Server) {
	error := server.ListenAndServe()
	if error != nil && error != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "[ERROR]: The server encountered an error: %s", error)
	}

	waitGroup.Done()
}

func interactiveConsole(server *http.Server) {
	scanner := bufio.NewScanner(os.Stdin)
	running := true

	for running {
		fmt.Printf("> ")
		scanner.Scan()
		command := scanner.Text()

		if command == "q" {
			running = false
			server.Shutdown(nil)
		}
	}
}

func main() {
	globalCardStacks = make(map[string]CardStack)

	serverMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":3000",
		Handler: serverMux,
	}

	serverMux.HandleFunc("/", staticFilesHandler)
	serverMux.HandleFunc("/stack/", stackPageHander)

	serverMux.HandleFunc("/api/cardstacks", apiCardStacksHandler)
	serverMux.HandleFunc("/api/cardstacks/", apiSpecificCardStackHandler)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	fmt.Printf("Starting server at http://127.0.0.1:3000\n")

	go runServer(&waitGroup, &server)
	interactiveConsole(&server)

	waitGroup.Wait()
}
