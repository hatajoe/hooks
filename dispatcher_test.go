package hooks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hatajoe/hooks/github"
)

func ExampleDispatcher_HTTPMux() {

	mux := http.NewServeMux()

	// Google AppEngine special endpoint
	mux.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Dispatcher is http.Handler
	// Dispatcher only routes under `/webhooks` pattern
	dispatcher := NewDispatcher(github.EventParser{})
	mux.Handle("/webhooks", dispatcher)

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

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

func TestDispatcher_ServeHTTP(t *testing.T) {

	mux := http.NewServeMux()

	// Google AppEngine special endpoint
	mux.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "/_ah/start")
	})

	// Dispatcher is http.Handler
	// Dispatcher only routes under `/webhooks` pattern
	dispatcher := NewDispatcher(github.EventParser{})
	dispatcher.On("echo", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
	})

	mux.Handle("/webhooks", dispatcher)

	httpReq := httptest.NewRequest("POST", "/webhooks", strings.NewReader("Hello"))
	httpReq.Header.Set("X-GitHub-Event", "echo")

	respWriter := stringResponseWriter{
		h: make(http.Header),
	}

	mux.ServeHTTP(&respWriter, httpReq)

	if respWriter.b.String() != "Hello" {
		t.Fatal(respWriter.b.String())
	}


	httpReq = httptest.NewRequest("POST", "/_ah/start", nil)
	respWriter.b.Reset()

	mux.ServeHTTP(&respWriter, httpReq)

	if respWriter.b.String() != "/_ah/start" {
		t.Fatal(respWriter.b.String())
	}
}
