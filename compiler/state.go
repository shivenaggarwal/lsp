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

// severity 1: error, 2: hint, 3: info, 4:warning
func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {

		// VS Code slander
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Please don't mention VS Code, no self-respecting dev uses it.",
			})
		}

		// Neovim appreciation
		if strings.Contains(line, "Neovim") {
			idx := strings.Index(line, "Neovim")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: 2,
				Source:   "Common Sense",
				Message:  "Great choice dawg.",
			})
		}

		// Anti-humility check
		if strings.Contains(line, "this file has no purpose") {
			idx := strings.Index(line, "this file has no purpose")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("this file has no purpose")),
				Severity: 3,
				Source:   "Self-Esteem Engine",
				Message:  "This file is *the* purpose. Believe in your LSP era.",
			})
		}

		// Vibe Check
		if strings.Contains(line, "not really") || strings.Contains(line, "yep. this is it") {
			idx := strings.Index(line, "not really")
			if idx < 0 {
				idx = strings.Index(line, "yep. this is it")
			}
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("not really")),
				Severity: 3,
				Source:   "VibeCheck",
				Message:  "Low-energy response detected. Inject more chaos.",
			})
		}

		// Slang Reminder
		if strings.Contains(line, "look, i made a thing") {
			idx := strings.Index(line, "look, i made a thing")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("look, i made a thing")),
				Severity: 3,
				Source:   "FlexPolice",
				Message:  "Say it with your chest. Try a flex format like: ‘caught in 4k making history.’",
			})
		}

		// Main character energy missing
		if strings.Contains(line, "congratulations") && strings.Contains(line, "underwhelming") {
			idx := strings.Index(line, "underwhelming")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("underwhelming")),
				Severity: 2,
				Source:   "Narrative Core",
				Message:  "Bro you’re underselling it. This is the climax of your origin story.",
			})
		}

		// No markdown heading enthusiasm
		if strings.HasPrefix(line, "# this is a test") {
			idx := strings.Index(line, "this is a test")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("this is a test")),
				Severity: 3,
				Source:   "Drama Department",
				Message:  "Let’s spice this up. How about: ‘Top-tier markdown drama incoming’?",
			})
		}
	}

	return diagnostics
}


func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	
	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
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

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		// social/internet-style completions
		{
			Label:         "caught in 4k absolutely feral frfr no cap in my main character era actually built different with receipts",
			Detail:        "Story template",
			Documentation: "Dramatic post format for attention-grabbing moments.",
		},
		{
			Label:         "can’t even lie that was diabolical literally in my regret era rn send help pls",
			Detail:        "Rant format",
			Documentation: "Unfiltered rant layout for expressive markdown.",
		},
		{
			Label:         "POV you just got absolutely obliterated by a side quest built different main quest energy",
			Detail:        "POV scenario",
			Documentation: "Turn anything into a cinematic experience.",
		},
		{
			Label:         "ratio’d into another dimension no cap that was savage af and I’m not recovering anytime soon",
			Detail:        "Ratio response",
			Documentation: "For spicy replies and comeback messages.",
		},
		{
			Label:         "flexing on the timeline like it’s fashion week living in my drip era no budget edition",
			Detail:        "Flex post",
			Documentation: "When you need to quietly brag with flair.",
		},
		{
			Label:         "tea spilled chaotic and undeniably iconic catch me refreshing every 3 seconds",
			Detail:        "Tea spill",
			Documentation: "Use for drama, gossip, or release notes leaks.",
		},
		{
			Label:         "you failed the vibe check unhinged energy detected please log off respectfully",
			Detail:        "Vibe check",
			Documentation: "Use to call out cursed or chaotic content.",
		},
		{
			Label:         "fit check initiated underrated slay in my closet cleanout arc serving lowkey couture",
			Detail:        "Fit check",
			Documentation: "Style callout and aesthetic review format.",
		},
		{
			Label:         "this user understood the assignment no notes ten outta ten certified slay moment",
			Detail:        "Slay report",
			Documentation: "Praise or hype up good execution.",
		},
		{
			Label:         "exposed thread undeniable evidence with screenshots time to get cancelled",
			Detail:        "Exposed thread",
			Documentation: "Ideal for receipts, bugs, or callouts.",
		},
		{
			Label:         "this take is spicy fully cooked and hotter than the sun report to opinion jail",
			Detail:        "Hot take",
			Documentation: "For strong opinions in markdown or commits.",
		},

		// professional-style completions
		{
			Label:         "devlog added feature X on component Y looks clean totally didn’t break anything yet",
			Detail:        "Developer log",
			Documentation: "Semi-serious changelog-style entry.",
		},
		{
			Label:         "TODO: implement feature before deadline prioritize accordingly",
			Detail:        "Task note",
			Documentation: "Mark tasks for future work cleanly.",
		},
		{
			Label:         "NOTE: potential issue with module under high load needs review",
			Detail:        "Code note",
			Documentation: "Surface potential problems with clarity.",
		},
		{
			Label:         "Pros: efficient, scalable, and easy to maintain",
			Detail:        "Positive analysis",
			Documentation: "Summarize benefits of an approach.",
		},
		{
			Label:         "Cons: brittle, verbose, and hard to debug",
			Detail:        "Negative analysis",
			Documentation: "List drawbacks honestly.",
		},
		{
			Label:         "Use Case: applicable when handling streaming data and requires low-latency processing",
			Detail:        "Usage scenario",
			Documentation: "Great for documentation or README context.",
		},
		{
			Label:         "Installation Instructions: run make install on Linux or use brew for MacOS",
			Detail:        "Setup guide",
			Documentation: "Quick install notes for onboarding.",
		},
		{
			Label:         "Example Usage: mycli --source ./data --output results.json",
			Detail:        "Code snippet",
			Documentation: "Embed in markdown for quick usage ref.",
		},
		{
			Label:         "Changelog: updated parser to v2.1 with better error recovery and more accurate output",
			Detail:        "Release note",
			Documentation: "For structured updates in projects.",
		},
		{
			Label:         "Summary: this document covers async design patterns including objectives and outcomes",
			Detail:        "Documentation header",
			Documentation: "Intro blurb for any markdown spec or design doc.",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
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
