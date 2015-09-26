package serfer

import (
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	tomb "gopkg.in/tomb.v2"
)

func TestRunSerfer(t *testing.T) {

	// Create event
	evt := &MockEvent{}

	// Create handler
	handler := &MockEventHandler{}
	handler.On("HandleEvent", evt).Return()

	// Create channel and serfer
	ch := make(chan serf.Event, 1)
	serfer := NewSerfer(ch, handler)

	// Setup test
	var death tomb.Tomb
	ctx, cancel := context.WithCancel(context.Background())

	// Start serfer
	death.Go(func() error {
		return serfer.Run(ctx)
	})

	// Send events
	select {
	case ch <- evt:
	case <-time.After(time.Second):
		t.Fatal("Event was not sent over channel")
	}
	ch <- evt

	// Stop event processing
	cancel()

	// Verify stopped without error
	assert.Nil(t, death.Wait(), "Error should be nil")

	// Validate event was prcoessed
	handler.AssertCalled(t, "HandleEvent", evt)

}

func TestRunSerfer_NilContext(t *testing.T) {

	// Create handler
	handler := &MockEventHandler{}

	// Create channel and serfer
	ch := make(chan serf.Event)
	serfer := NewSerfer(ch, handler)

	// Verify stopped with error
	assert.NotNil(t, serfer.Run(nil), "Error should not be nil")
	handler.AssertNotCalled(t, "HandleEvent")
}
