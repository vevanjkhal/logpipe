package pipeline

import (
	"github.com/yourorg/logpipe/internal/routing"
)

// WithRouting attaches a Router to the Builder so that lines are
// dispatched to different outputs based on content.
// Call RegisterRouteSink after this to map destination labels to writers.
func (b *Builder) WithRouting(router *routing.Router) *Builder {
	b.router = router
	return b
}

// RegisterRouteSink maps a destination label to an already-added writer
// index. Panics if WithRouting has not been called first.
func (b *Builder) RegisterRouteSink(destination string, w interface{ Write([]byte) (int, error) }) *Builder {
	if b.router == nil {
		panic("pipeline: RegisterRouteSink called before WithRouting")
	}
	if b.dispatcher == nil {
		b.dispatcher = routing.NewDispatcher(b.router, b.fallbackWriter())
	}
	b.dispatcher.Register(destination, w)
	return b
}
