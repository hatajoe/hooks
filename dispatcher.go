package hooks

import (
	"fmt"
	"net/http"
)

type dispatcher struct {
	eventParser func(r *http.Request) (string, error)
	handlers    map[string]http.HandlerFunc
}

// NewDispatcher returns the event dispatcher object
// The argument `eventParser` is that parses event string from *http.Request
func NewDispatcher(eventParser func(r *http.Request) (string, error)) *dispatcher {
	return &dispatcher{
		eventParser: eventParser,
		handlers:    map[string]http.HandlerFunc{},
	}
}

// On adds a handler corresponding the specific event string
func (d *dispatcher) On(event string, handler http.HandlerFunc) {
	d.handlers[event] = handler
}

// Listen starts HTTP server that handling the registered event
// The first argument `endpoint` is the path of the hooks URI (e.g, "/webhooks")
// The second argument `port` is a listen port (e.g, ":3000")
func (d dispatcher) Listen(endpoint, port string) error {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		e, err := d.eventParser(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if _, ok := d.handlers[e]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("event is not registered: `%s`", e)))
			return
		}
		d.handlers[e](w, r)
	})
	return http.ListenAndServe(port, nil)
}
