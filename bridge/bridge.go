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
	tx_cmds       chan *Command
	handlers      []chan *Command
	handlers_lock *sync.RWMutex
}

func New() *Bridge {
	return &Bridge{
		tx_cmds:       make(chan *Command, 1000),
		handlers:      make([]chan *Command, 0),
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
	b.tx_cmds <- &Command{Type: t, Data: data}
	return nil
}

// execute a command on the remote side and wait for the response
func (b *Bridge) ExecuteCommand(t CommandType, obj interface{}, timeout time.Duration) (*Command, error) {
	var rx_cmd *Command
	var err error
	id := math.Floor(rand.Float64() * 64000)

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	b.tx_cmds <- &Command{Type: t, Data: data, ID: &id}

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

type CommandHandler func(*Command)

func (b *Bridge) AddHandler() chan *Command {
	h := make(chan *Command, 100)
	b.handlers_lock.Lock()
	b.handlers = append(b.handlers, h)
	b.handlers_lock.Unlock()
	return h
}

func (b *Bridge) RemoveHandler(remove_handler chan *Command) {
	newhandlers := make([]chan *Command, 0)
	b.handlers_lock.Lock()
	for _, h := range b.handlers {
		if h != remove_handler {
			newhandlers = append(newhandlers, h)
		}
	}
	b.handlers = newhandlers
	b.handlers_lock.Unlock()
}
