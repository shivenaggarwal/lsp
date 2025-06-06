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

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

// this is what our lsp will reply with
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	// TextDocumentSync int `json:"textDocumentSync"` // TODO incremental updates, etc
	TextDocumentSync TextDocumentSyncOptions `json:"textDocumentSync"`

	HoverProvider bool `json:"hoverProvider"`
	DefinitionProvider bool `json:"definitionProvider"`
	CodeActionProvider bool `json:"codeActionProvider"` //TODO
}

type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose"`
	Change    int  `json:"change"` // 1 = Full, 2 = Incremental
	Save      SaveOptions `json:"save"`
}

type SaveOptions struct {
	IncludeText bool `json:"includeText"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
			
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    1, // Full document sync
					Save: SaveOptions{
						IncludeText: true,
					},
				},
				HoverProvider: true,
				DefinitionProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "lsp",
				Version: "0.0.1-beta1.final",
			},
		},
	}
}




