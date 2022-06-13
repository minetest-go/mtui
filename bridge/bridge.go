package bridge

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Bridge struct {
	tx_cmds       chan *CommandRequest
	handlers      []chan *CommandResponse
	handlers_lock *sync.RWMutex
}

func New() *Bridge {
	return &Bridge{
		tx_cmds:       make(chan *CommandRequest, 1000),
		handlers:      make([]chan *CommandResponse, 0),
		handlers_lock: &sync.RWMutex{},
	}
}

var ErrTimeout = errors.New("timeout")

// one-way command, no response
func (b *Bridge) SendCommand(t CommandRequestType, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	b.tx_cmds <- &CommandRequest{Type: t, Data: data}
	return nil
}

// execute a command on the remote side and wait for the response
func (b *Bridge) ExecuteCommand(t CommandRequestType, obj interface{}, timeout time.Duration) (*CommandResponse, error) {
	var rx_cmd *CommandResponse
	var err error
	id := math.Floor(rand.Float64() * 64000)

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	b.tx_cmds <- &CommandRequest{Type: t, Data: data, ID: &id}

	c := b.AddHandler()
	then := time.Now().Add(timeout)

	for {
		select {
		case cmd := <-c:
			if cmd.ID != nil && *cmd.ID == id {
				rx_cmd = cmd
				break
			}
		case <-time.After(100 * time.Millisecond):
		}
		if rx_cmd != nil {
			// result received
			break
		}
		if time.Now().After(then) {
			// request timed out
			err = ErrTimeout
			break
		}
	}

	b.RemoveHandler(c)
	return rx_cmd, err
}

type CommandHandler func(*CommandResponse)

func (b *Bridge) AddHandler() chan *CommandResponse {
	h := make(chan *CommandResponse, 100)
	b.handlers_lock.Lock()
	b.handlers = append(b.handlers, h)
	b.handlers_lock.Unlock()
	return h
}

func (b *Bridge) RemoveHandler(remove_handler chan *CommandResponse) {
	newhandlers := make([]chan *CommandResponse, 0)
	b.handlers_lock.Lock()
	for _, h := range b.handlers {
		if h != remove_handler {
			newhandlers = append(newhandlers, h)
		}
	}
	b.handlers = newhandlers
	b.handlers_lock.Unlock()
}
