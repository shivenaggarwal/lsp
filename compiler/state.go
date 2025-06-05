package compiler

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
