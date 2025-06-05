package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`

	// TODO add more params to handle a real language
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
