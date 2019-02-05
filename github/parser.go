package github

import (
	"fmt"
	"net/http"
)

// Parser parse GitHub event from http.Request
// Parser is `EventParser` implementation
type Parser struct{}

// GetEvent returns the event string from http.Request
func (Parser) GetEvent(r *http.Request) (string, error) {
	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		return "", fmt.Errorf("X-GitHub-Event is empty")
	}
	return event, nil
}
