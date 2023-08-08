package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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

func main() {
	globalCardStacks = make(map[string]CardStack)

	http.HandleFunc("/", staticFiles)
	http.HandleFunc("/stack/", stackPage)

	http.HandleFunc("/api/cardstacks", apiCardStacks)
	http.HandleFunc("/api/cardstacks/", apiSpecificCardStack)

	error := http.ListenAndServe(":3000", nil)
	if error != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: The server encountered an error: %s", error)
	}
}
