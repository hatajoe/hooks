package github

import (
	"net/http/httptest"
	"testing"
)

func TestParserCorrect(t *testing.T) {

	expect := "push"

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-GitHub-Event", expect)

	p := &EventParser{}
	if actual, err := p.Parse(req); err != nil {
		t.Errorf("%v", err)
	} else if expect != actual {
		t.Errorf("actual value is not expected: `%s`", actual)
	}
}

func TestParserFailure(t *testing.T) {
	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Event", "push")

	p := &EventParser{}
	if _, err := p.Parse(req); err == nil {
		t.Errorf("Error occurring is expected, but no error was detected.")
	}
}
