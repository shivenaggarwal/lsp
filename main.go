package main

import (
	"bufio"
	"encoding/json"
	"log"
	"lsp/lsp"
	"lsp/rpc"
	"os"
)

func main() {
	logger := getLogger("/Users/shiven/Developer/lsp/log.txt")
	logger.Println("Logger Started!")

	scanner := bufio.NewScanner(os.Stdin) // reading input from stdin
	scanner.Split(rpc.Split)              // now this is going to be waiting on stdout for the new msg and once we get that the scanner will give use the text and then we can handle it below

	for scanner.Scan() { // basically saying keep on running it until we are ready
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Panicf("got an error: %s", err)
			continue
		}
		HandleMessage(logger, method, contents)
	}
}

func HandleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("received message with method: %s", method) // this just makes sure everytime we get a message we print the message, lets us know are we decoding msg and passing them forward correctly
	
	switch method {
		case "initialize": 
			var request lsp.InitializeRequest
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("hey, couldn't parse this: %s", err)
			}
			logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

			// once we know we have a message i.e a request our lsp should reply
			msg := lsp.NewInitializeResponse(request.ID)
			reply := rpc.EncodeMessage(msg)
			writer := os.Stdout
			writer.Write([]byte(reply)) // this is going to take the stdout of the current process and reply back with this a sequence of bytes
			logger.Print("sent the reply")
	}

}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666) // just making it so it is open to all users and it has read and write both

	if err != nil {
		panic("bad file, give good file")
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
