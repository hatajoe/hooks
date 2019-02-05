package hooks

import (
	"fmt"
	"net/http"
)

// EventParser is the event parser interface of http.Request
type EventParser interface {
	// Parse parses the event name from http.Request
	Parse(r *http.Request) (string, error)
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
		e, err := d.eventParser.Parse(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if handler, ok := d.handlers[e]; !ok {
			http.Error(w, fmt.Sprintf("event is not registered: `%s`", e), http.StatusBadRequest)
			return
		} else {
			handler(w, r)
		}
	})
	return http.ListenAndServe(port, nil)
}
