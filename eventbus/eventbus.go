package eventbus

import (
	"sync"
)

type EventType string

type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

type EventBus struct {
	mutex     *sync.RWMutex
	listeners []chan *Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		mutex:     &sync.RWMutex{},
		listeners: make([]chan *Event, 0),
	}
}

func (eb *EventBus) AddListener(ch chan *Event) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.listeners = append(eb.listeners, ch)
}

func (eb *EventBus) RemoveListener(ch chan *Event) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	new_listeners := make([]chan *Event, 0)
	for _, existing_listener := range eb.listeners {
		if existing_listener != ch {
			new_listeners = append(new_listeners, existing_listener)
		}
	}

	eb.listeners = new_listeners
}

func (eb *EventBus) Emit(wse *Event) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	for _, ch := range eb.listeners {
		select {
		case ch <- wse:
		default:
		}
	}
}
