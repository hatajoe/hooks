package slack

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type stringResponseWriter struct {
	h http.Header
	b strings.Builder
	c int
}

func (w *stringResponseWriter) Header() http.Header {
	return w.h
}

func (w *stringResponseWriter) Write(b []byte) (int, error) {
	return w.b.Write(b)
}

func (w *stringResponseWriter) WriteHeader(code int) {
	w.c = code
}

func TestChallengeHandler(t *testing.T) {
	f, err := os.Open("./fixtures/challenge.json")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/webhooks", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
	rec := &stringResponseWriter{
		h: make(http.Header),
	}

	ChallengeHandler(rec, req)

	if rec.c != http.StatusOK {
		t.Errorf("Unexpected response code: %d", rec.c)
	}

	if rec.b.String() != "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P" {
		t.Errorf("Unexpected challenge response: %s", rec.b.String())
	}
}
