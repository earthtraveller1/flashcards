package main

import (
	"bufio"
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

type CardStack struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cards       []Card `json:"cards"`
}

var globalCardStacks map[string]CardStack

func staticFiles(pWriter http.ResponseWriter, pRequest *http.Request) {
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

func stackPage(pWriter http.ResponseWriter, pRequest *http.Request) {
	uri := pRequest.URL.RequestURI()
	uriParts := strings.Split(uri, "/")
	stackName := uriParts[len(uriParts)-1]

	_, stackExists := globalCardStacks[stackName]
	if !stackExists {
		pWriter.WriteHeader(404)
		fmt.Fprintf(pWriter, "<h1>404</h1> <p>Not found")
		return
	}

	stackPageTemplateBytes, error := os.ReadFile("stack.html")
	if error != nil {
		pWriter.WriteHeader(500)
		pWriter.Write([]byte("<h1>500</h1> <p>Internal server error"))

		fmt.Fprintf(os.Stderr, "[ERROR]: stack.html does not appear to exist, or I cannot load it or something.\n")
	}

	stackPageTemplate := string(stackPageTemplateBytes)

	substitution := "//{{{server}}}"
	serverParams := `
const serverInfo = {
    stackName: "%s"
}
    `

	serverParams = fmt.Sprintf(serverParams, stackName)

	var stackPage string
	if strings.Contains(stackPageTemplate, substitution) {
		stackPage = strings.Replace(stackPageTemplate, substitution, serverParams, 1)
	}

	pWriter.Write([]byte(stackPage))
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

		if strings.HasPrefix(command, "shutdown") {
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

	serverMux.HandleFunc("/", staticFiles)
	serverMux.HandleFunc("/stack/", stackPage)

	serverMux.HandleFunc("/api/cardstacks", apiCardStacks)
	serverMux.HandleFunc("/api/cardstacks/", apiSpecificCardStack)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	go runServer(&waitGroup, &server)
	interactiveConsole(&server)

	waitGroup.Wait()
}
