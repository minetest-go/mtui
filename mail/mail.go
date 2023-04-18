package mail

import (
	"encoding/json"
	"fmt"

	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/mod_storage"
)

type Mail struct {
	ctx *mtdb.Context
}

func New(ctx *mtdb.Context) *Mail {
	return &Mail{ctx: ctx}
}

func (m *Mail) GetEntry(playername string) (*PlayerEntry, error) {
	e, err := m.ctx.ModStorage.Get("mail", []byte(fmt.Sprintf("mail/%s", playername)))
	pe := &PlayerEntry{
		Inbox:    make([]*Message, 0),
		Outbox:   make([]*Message, 0),
		Drafts:   make([]*Message, 0),
		Contacts: make([]*Contact, 0),
		Lists:    make([]*Maillist, 0),
	}
	if e == nil {
		// return empty entry
		return pe, nil
	}
	if err != nil {
		return nil, err
	}
	return pe, json.Unmarshal(e.Value, pe)
}

func (m *Mail) SetEntry(playername string, pe *PlayerEntry) error {
	data, err := json.Marshal(pe)
	if err != nil {
		return err
	}

	e, err := m.ctx.ModStorage.Get("mail", []byte(fmt.Sprintf("mail/%s", playername)))
	if err != nil {
		return err
	}

	if e == nil {
		// create new
		return m.ctx.ModStorage.Create(&mod_storage.ModStorageEntry{
			ModName: "mail",
			Key:     []byte(fmt.Sprintf("mail/%s", playername)),
			Value:   data,
		})
	} else {
		// update existing
		e.Value = data
		return m.ctx.ModStorage.Update(e)
	}
}
