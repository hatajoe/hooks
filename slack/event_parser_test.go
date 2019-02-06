package slack

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParser_Parse_Correct(t *testing.T) {

	expect := "app_mention"

	f, err := os.Open("./fixtures/app_mention.json")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	p := &EventParser{
		ChallengeEventType: "challenge",
		VerifyToken:        true,
		VerificationToken:  "test-token",
	}

	if actual, err := p.Parse(req); err != nil {
		t.Errorf("%v", err)
	} else if expect != actual {
		t.Errorf("actual value is not expected: `%s`", actual)
	}
}

func TestParser_Parse_Failure(t *testing.T) {
	f, err := os.Open("./fixtures/app_mention.json")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	p := &EventParser{
		ChallengeEventType: "challenge",
		VerifyToken:        true,
		VerificationToken:  "testtoken",
	}

	if _, err := p.Parse(req); err == nil {
		t.Errorf("returning error is expected, but no error occurred")
	}
}
