package bridge

import (
	"mtui/types"
	"sync"
)

type Bridge struct {
	tx_cmds       chan *types.Command
	handlers      []CommandHandler
	handlers_lock *sync.RWMutex
}

func New() *Bridge {
	return &Bridge{
		tx_cmds:       make(chan *types.Command, 1000),
		handlers:      make([]CommandHandler, 0),
		handlers_lock: &sync.RWMutex{},
	}
}

func (b *Bridge) SendCommand(cmd *types.Command) {
	b.tx_cmds <- cmd
}

type CommandHandler func(*types.Command)

func (b *Bridge) RegisterCommandHandler(handler CommandHandler) {
	b.handlers_lock.Lock()
	b.handlers = append(b.handlers, handler)
	b.handlers_lock.Unlock()
}
