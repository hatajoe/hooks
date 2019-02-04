package hooks

import (
	"fmt"
	"io"
	"net/http"
)

// EventParser is the event parser interface of http.Request
type EventParser interface {
	GetEvent(r *http.Request) (string, error)
}

// Dispatcher is HTTP server that handles the event of http.Request
type Dispatcher struct {
	eventParser EventParser
	handlers    map[string]http.HandlerFunc
}

// NewDispatcher returns the event Dispatcher object
// The argument `p` is that parses event string from *http.Request
func NewDispatcher(p EventParser) *Dispatcher {
	return &Dispatcher{
		eventParser: p,
		handlers:    map[string]http.HandlerFunc{},
	}
}

// On adds a handler corresponding the specific event string
func (d *Dispatcher) On(event string, handler http.HandlerFunc) {
	d.handlers[event] = handler
}

// Listen starts HTTP server that handling the registered event
// The first argument `endpoint` is the path of the hooks URI (e.g, "/webhooks")
// The second argument `port` is a listen port (e.g, ":3000")
func (d Dispatcher) Listen(endpoint, port string) error {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		e, err := d.eventParser.GetEvent(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fprintf(w, "%s", err.Error())
			return
		}
		if handler, ok := d.handlers[e]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fprintf(w, "event is not registered: `%s`", e)
			return
		} else {
			handler(w, r)
		}
	})
	return http.ListenAndServe(port, nil)
}

// fprintf is ignore return value of fmt.Fprintf
func fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}
