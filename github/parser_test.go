package github

import (
	"net/http/httptest"
	"testing"
)

func TestParserCorrect(t *testing.T) {

	expect := "push"

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-GitHub-Event", expect)

	p := &Parser{}
	if actual, err := p.GetEvent(req); err != nil {
		t.Errorf("%v", err)
	} else if expect != actual {
		t.Errorf("actual value is not expected: `%s`", actual)
	}
}

func TestParserFailure(t *testing.T) {
	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Header.Set("X-Event", "push")

	p := &Parser{}
	if _, err := p.GetEvent(req); err == nil {
		t.Errorf("Error occurring is expected, but no error was detected.")
	}
}
