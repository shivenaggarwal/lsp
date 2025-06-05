package lsp

// we send a response for each request type 
type Request struct {
	RPC string `json:"jsonrpc"`
	ID int `json:"id"`
	Method string `json:"method"`

	//TODO specify the params in all request types
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID *int `json:"id,omitempty"` // id could be empty here

	// some result TODO
	// some error TODO
}

type Notification struct {
	RPC string `json:"jsonrpc"`
	Method string `json:"method"`
}
