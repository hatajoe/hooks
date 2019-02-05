package github

import (
	"fmt"
	"net/http"
)

// EventParser parse GitHub event from http.Request
// EventParser is `EventParser` implementation
type EventParser struct{}

// Parse returns the event string from http.Request
func (EventParser) Parse(r *http.Request) (string, error) {
	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		return "", fmt.Errorf("X-GitHub-Event is empty")
	}
	return event, nil
}
