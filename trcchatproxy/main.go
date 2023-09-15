package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/trimble-oss/tierceron-nute/mashupsdk"
	"github.com/trimble-oss/tierceron-nute/mashupsdk/server"
)

//go:embed tls/mashup.crt
var mashupCert embed.FS

//go:embed tls/mashup.key
var mashupKey embed.FS

var gchatApp GChatApp
var id int64

func (w *GChatApp) InitServer(callerCreds string, insecure bool, maxMessageLength int) {
	if callerCreds != "" {
		server.InitServer(callerCreds, insecure, maxMessageLength, w.MashupSdkApiHandler, w.WClientInitHandler)
	}
}

func main() {
	gchatApp = GChatApp{
		MashupSdkApiHandler: &GoogleChatHandler{},
		GoogleChatContext:   &GoogleChatContext{},
		WClientInitHandler:  &WorldClientInitHandler{},
	}
	shutdown := make(chan bool)

	// Initialize local server.
	mashupsdk.InitCertKeyPair(mashupCert, mashupKey)

	gchatworld := GChatApp{
		MashupSdkApiHandler: &GoogleChatHandler{},
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	configPort, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	configs := mashupsdk.MashupConnectionConfigs{
		AuthToken:   "",
		CallerToken: "",
		Server:      "",
		Port:        configPort,
	}
	encoding, err := json.Marshal(&configs)
	if err != nil {
		fmt.Println(err)
	}

	callerCreds := flag.String("CREDS", string(encoding), "Credentials of caller")
	flag.Parse()
	id = 0
	server.RemoteInitServer(*callerCreds, true, -2, gchatworld.MashupSdkApiHandler, gchatworld.WClientInitHandler)

	<-shutdown
}

// Processes upserted query from client
// Changes based on msg.Name
func ProcessQuery(msg *mashupsdk.MashupDetailedElement) {
	if msg.Name == "DialogFlow" {
		ProcessDFQuery(msg)
	} else if msg.Name == "DialogFlowResponse" {
		ProcessDFResponse(msg)
	} else if msg.Name == "GChatResponse" {
		ProcessGChatAnswer(msg)
	} else if msg.Name == "Get Message" {
		gchatApp.DetailedElements = gchatApp.DetailedElements[:len(gchatApp.DetailedElements)-1]
		input := ""
		for input == "" {
			input = getUserInput()
			if input != "" {
				gchatApp.DetailedElements = append(gchatApp.DetailedElements, &mashupsdk.MashupDetailedElement{
					Name: "GChatQuery",
					Id:   int64(len(gchatApp.DetailedElements)), // Make sure id matches index in elements
					Data: input,
				})
			} else {
				fmt.Println("An error occurred with reading the input. Please input your question in the command line and press enter!")
			}
		}
	} else {
		log.Printf("Message type does not correspond to either GChatQuery or DialogFlow")
	}
}

// Asks user input
// This is a stub version --> potentially shouldn't be needed if user can @askflume in google chat
// However, maybe use this as a way to ask user if there is anything else they would like to ask
func getUserInput() string {
	fmt.Println("This is a simulation of the Flume Chat App. Please type your question below and press enter: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input from user: %v", err)
		return ""
	}
	return input
}

// Updates ID and returns value
// id should match up with number of queries made by user
func GetId() int64 {
	id += 1
	return id
}