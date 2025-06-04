package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// takes in some message type and returns a string
func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg) // used to convert msg to json format
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

// every message we send or receive will have a method key inside of the json struct, we index off of that to decide what else to do
type BaseMessage struct {
	Method string `json:"method"`
}

// take the json format and converts into an actual value
func DecodeMessage(msg []byte) (string, []byte, error) {
	// bytes.cut takes the msg, cut slices around the first seperator and gives back all the byts before and after that and a check bool which tells whether we found it or not
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found{
		return "", nil, errors.New("Did not find seperator")
	}
		
	// Content-Length: <number> 
	contentLengthBytes := header[len("Content-Length: "):] // gives back the bytes after content length header
	contentLength, err := strconv.Atoi(string(contentLengthBytes)) // converting bytes string to integer
	if err != nil {
		return "", nil, err
	}	

	var baseMessage BaseMessage
	if err:= json.Unmarshal(content[:contentLength], &baseMessage); err != nil { // unpacking the message with unmarshall
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}
