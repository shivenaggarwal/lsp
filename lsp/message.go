package lsp

type Request struct {
	RPC string `json:"jsonrpc"`
	ID int `json:"id"`
	Method string `json:"method"`

	//TODO specify the params in all request types
}
