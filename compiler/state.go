package compiler

import (
	"lsp/lsp"
	"fmt"
	"strings"
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

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// this would look up the position
	// but for now im gonna put this as definition to be one line above

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]
	lines := strings.Split(text, "\n")

	actions := []lsp.CodeAction{}

	for row, line := range lines[:len(lines)-1] { // exclude the last line
		// Replace "VS Code" with "Neovim"
		if idx := strings.Index(line, "VS Code"); idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{
				uri: {{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				}},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{
				uri: {{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				}},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}

		// Underwhelming → masterpiece
		if idx := strings.Index(line, "underwhelming"); idx >= 0 {
			actions = append(actions, lsp.CodeAction{
				Title: "Replace 'underwhelming' with 'masterpiece'",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, idx, idx+len("underwhelming")),
						NewText: "masterpiece",
					}},
				}},
			})
		}

		// Answer rhetorical question
		if strings.Contains(line, "does it do anything cool") {
			actions = append(actions, lsp.CodeAction{
				Title: "Answer rhetorical question with sarcasm",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row+1, 0, 0),
						NewText: "\n> define 'cool'.",
					}},
				}},
			})
		}

		// Strikethrough self-deprecation
		if strings.Contains(line, "this file has no purpose") {
			actions = append(actions, lsp.CodeAction{
				Title: "Strike through self-deprecating line",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, 0, len(line)),
						NewText: "~~" + line + "~~",
					}},
				}},
			})
		}

		// Add dramatic emoji to ends with "." or "..."
		if strings.HasSuffix(line, "...") || strings.HasSuffix(line, ".") {
			actions = append(actions, lsp.CodeAction{
				Title: "Add dramatic emoji for flair",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, len(line), len(line)),
						NewText: " 👀",
					}},
				}},
			})
		}

		// Replace markdown header with emoji title
		if strings.HasPrefix(line, "# this is a test") {
			actions = append(actions, lsp.CodeAction{
				Title: "Replace header with something dramatic",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, 0, len(line)),
						NewText: "# the README that read too much into itself",
					}},
				}},
			})
		}

		// Emphasize ego line
		if strings.Contains(line, "boosting my ego") {
			actions = append(actions, lsp.CodeAction{
				Title: "Emphasize ego line with blockquote",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, 0, len(line)),
						NewText: "> " + line,
					}},
				}},
			})
		}

		// Add fake TODO after the title
		if row == 1 {
			actions = append(actions, lsp.CodeAction{
				Title: "Add fake TODO to look productive",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row+1, 0, 0),
						NewText: "\n<!-- TODO: add real content to this file -->\n\n",
					}},
				}},
			})
		}

		// Replace "yep. this is it." with dramatic message
		if strings.Contains(line, "yep. this is it") {
			actions = append(actions, lsp.CodeAction{
				Title: "Replace 'yep. this is it.' with something dramatic",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, 0, len(line)),
						NewText: "brace yourself. markdown greatness incoming.",
					}},
				}},
			})
		}

		// Italicize "cool"
		if idx := strings.Index(line, "cool"); idx >= 0 {
			actions = append(actions, lsp.CodeAction{
				Title: "Italicize 'cool' for ironic tone",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, idx, idx+len("cool")),
						NewText: "*cool*",
					}},
				}},
			})
		}

		// Replace "say 'look...'" line with bolder message
		if strings.Contains(line, "say \"look, i made a thing\"") {
			actions = append(actions, lsp.CodeAction{
				Title: "Make statement bolder",
				Edit: &lsp.WorkspaceEdit{Changes: map[string][]lsp.TextEdit{
					uri: {{
						Range:   LineRange(row, 0, len(line)),
						NewText: "brag: i built something. deal with it.",
					}},
				}},
			})
		}
	}

	return lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}


func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
