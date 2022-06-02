package bridge

import (
	"encoding/json"
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
	b.SendCommand(t, obj)
	return nil, nil //TODO
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
