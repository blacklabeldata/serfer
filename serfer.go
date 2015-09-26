package serfer

import (
	"errors"

	"github.com/hashicorp/serf/serf"
	"golang.org/x/net/context"
)

// Serfer processes Serf.Events and is meant to be ran in a goroutine.
type Serfer interface {

	// Run continuously reads from a channel of serf.Events. When the context is closed, the method should stop processing and return.
	Run(ctx context.Context) error
}

// NewSerfer returns a new Serfer implementation that uses the given channel and event handlers.
func NewSerfer(c chan serf.Event, handler EventHandler) Serfer {
	return &serfer{handler, c}
}

type serfer struct {
	handler EventHandler
	channel chan serf.Event
}

func (s *serfer) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context cannot be nil")
	}

	// Start event processing
	for {
		select {

		// Handle context close
		case <-ctx.Done():
			return nil

		// Handle serf events
		case evt := <-s.channel:
			s.handler.HandleEvent(evt)
		}
	}
}
