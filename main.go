package main

import (
	"bufio"
	"log"
	"lsp/rpc"
	"os"
)

func main()  {
	logger := getLogger("/Users/shiven/Developer/lsp/log.txt")
	logger.Println("Logger Started!")

	scanner:= bufio.NewScanner(os.Stdin) // reading input from stdin
	scanner.Split(rpc.Split) // now this is going to be waiting on stdout for the new msg and once we get that the scanner will give use the text and then we can handle it below

	for scanner.Scan() { // basically saying keep on running it until we are ready
		msg := scanner.Text()
		HandleMessage(logger, msg)
	}
}

func HandleMessage(logger *log.Logger, msg any) {
	logger.Println(msg) // this just makes sure everytime we get a message we print the message, lets us know are we decoding msg and passing them forward correctly	
}

func getLogger (filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666) // just making it so it is open to all users and it has read and write both 

	if err != nil {
		panic("bad file, give good file")
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
