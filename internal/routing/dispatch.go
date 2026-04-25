package routing

import "io"

// Dispatcher writes each log line to the output selected by a Router.
// Destinations that are not registered fall back to the default writer.
type Dispatcher struct {
	router   *Router
	sinks    map[string]io.Writer
	fallback io.Writer
}

// NewDispatcher creates a Dispatcher backed by router.
// fallback receives lines whose destination has no registered sink.
func NewDispatcher(router *Router, fallback io.Writer) *Dispatcher {
	return &Dispatcher{
		router:   router,
		sinks:    make(map[string]io.Writer),
		fallback: fallback,
	}
}

// Register associates a destination label with a writer.
func (d *Dispatcher) Register(destination string, w io.Writer) {
	d.sinks[destination] = w
}

// Dispatch routes line to the appropriate writer and writes it followed
// by a newline. It returns any write error encountered.
func (d *Dispatcher) Dispatch(line string) error {
	dest := d.router.Route(line)
	w, ok := d.sinks[dest]
	if !ok {
		w = d.fallback
	}
	_, err := io.WriteString(w, line+"\n")
	return err
}
