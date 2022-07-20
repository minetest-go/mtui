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
	handlers      map[CommandType][]chan *CommandResponse
	handlers_lock *sync.RWMutex
}

func New() *Bridge {
	return &Bridge{
		tx_cmds:       make(chan *CommandRequest, 1000),
		handlers:      make(map[CommandType][]chan *CommandResponse),
		handlers_lock: &sync.RWMutex{},
	}
}

var ErrTimeout = errors.New("timeout")

// one-way command, no response
func (b *Bridge) SendCommand(t CommandType, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	b.tx_cmds <- &CommandRequest{Type: t, Data: data}
	return nil
}

// execute a command on the remote side and wait for the response
func (b *Bridge) ExecuteCommand(t CommandType, obj any, resp any, timeout time.Duration) error {
	var rx_cmd *CommandResponse
	var err error
	id := math.Floor(rand.Float64() * 64000)

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	b.tx_cmds <- &CommandRequest{Type: t, Data: data, ID: &id}

	c := b.AddHandler(t)
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

	b.RemoveHandler(t, c)
	if err != nil {
		return err
	}
	return json.Unmarshal(rx_cmd.Data, resp)
}

type CommandHandler func(*CommandResponse)

func (b *Bridge) AddHandler(ct CommandType) chan *CommandResponse {
	b.handlers_lock.Lock()
	defer b.handlers_lock.Unlock()

	h := make(chan *CommandResponse, 100)
	l := b.handlers[ct]
	if l == nil {
		l = make([]chan *CommandResponse, 0)
	}
	l = append(l, h)
	b.handlers[ct] = l

	return h
}

func (b *Bridge) RemoveHandler(ct CommandType, remove_handler chan *CommandResponse) {
	b.handlers_lock.Lock()
	defer b.handlers_lock.Unlock()

	oldhandlers := b.handlers[ct]
	newhandlers := make([]chan *CommandResponse, 0)
	for _, h := range oldhandlers {
		if h != remove_handler {
			newhandlers = append(newhandlers, h)
		}
	}
	b.handlers[ct] = newhandlers
}
