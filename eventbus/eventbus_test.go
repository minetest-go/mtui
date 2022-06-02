package eventbus_test

import (
	"mtui/eventbus"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestEventType eventbus.EventType = "testevent"

func TestEventBus(t *testing.T) {
	e := eventbus.NewEventBus()

	ch := make(chan *eventbus.Event, 100)
	e.AddListener(ch)

	e.Emit(&eventbus.Event{
		Type: TestEventType,
		Data: 123,
	})

	evt := <-ch
	assert.NotNil(t, evt)
	assert.Equal(t, TestEventType, evt.Type)
	assert.Equal(t, 123, evt.Data)

	e.RemoveListener(ch)
}
