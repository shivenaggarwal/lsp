package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"lsp/compiler"
	"lsp/lsp"
	"lsp/rpc"
	"os"
)

func main() {
	logger := getLogger("/Users/shiven/Developer/lsp/log.txt")
	logger.Println("Logger Started!")

	scanner := bufio.NewScanner(os.Stdin) // reading input from stdin
	scanner.Split(rpc.Split)              // now this is going to be waiting on stdout for the new msg and once we get that the scanner will give use the text and then we can handle it below

	state := compiler.NewState()	
	writer := os.Stdout

	for scanner.Scan() { // basically saying keep on running it until we are ready
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Panicf("got an error: %s", err)
			continue
		}
		HandleMessage(logger, writer, state, method, contents)
	}
}

func HandleMessage(logger *log.Logger, writer io.Writer, state compiler.State, method string, contents []byte) {
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

			writeResponse(writer, msg)
			
			// debugging
			// logger.Printf("initialize response: %s", string(reply))	
			
			logger.Print("sent the reply")

		case "textDocument/didOpen":
			var request lsp.DidOpenTextDocumentNotification
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
			}
			// logger.Printf("Opened: %s %s", request.Params.TextDocument.URI, request.Params.TextDocument.Text)
			logger.Printf("Opened: %s", request.Params.TextDocument.URI)
			
			diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		
		case "textDocument/didChange":
			var request lsp.TextDocumentDidChangeNotification
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
			}
			logger.Printf("Changed: %s", request.Params.TextDocument.URI)
			for _, change := range request.Params.ContentChanges{
				
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}

		case "textDocument/hover":
			var request lsp.HoverRequest
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
			}
			
			// create a response and write it back
			response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position) // creating a response

			writeResponse(writer, response) // writing it back

		case "textDocument/definition":
			var request lsp.DefinitionRequest
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
			}
			
			// create a response and write it back
			response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position) // creating a response

			writeResponse(writer, response) // writing it back
			
		case "textDocument/codeAction":
			var request lsp.DefinitionRequest
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return
			}
			
			// create a response and write it back
			response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI) // creating a response

			writeResponse(writer, response) // writing it back
		
		case "textDocument/completion":
			var request lsp.CompletionRequest
			if err:= json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/completion: %s", err)
			return
			}
			
			// create a response and write it back
			response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI) // creating a response

			writeResponse(writer, response) // writing it back
			// debugging
			// logger.Printf("raw contents: %s", string(contents))
		
	}
}



func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply)) // this is going to take the stdout of the current process and reply back with this a sequence of bytes
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666) // just making it so it is open to all users and it has read and write both

	if err != nil {
		panic("bad file, give good file")
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
