package compiler

import (
	"lsp/lsp"
	"fmt"
	)

// basically to track the state of whats going on eg keeping track of our docs
type State struct {
	// map of filenames to content
	Documents map[string]string // whatever the current state of all of the current docs is we save them here
}

func NewState() State{
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// the function of this is to look up the type in our compiler code (will add it in future)

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}
