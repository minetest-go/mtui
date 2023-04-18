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

type PlayerEntry struct {
	Inbox    []*Message  `json:"inbox"`
	Outbox   []*Message  `json:"outbox"`
	Drafts   []*Message  `json:"drafts"`
	Contacts []*Contact  `json:"contacts"`
	Lists    []*Maillist `json:"lists"`
}

type Maillist struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}

type Message struct {
	ID      string  `json:"id"`
	Body    string  `json:"body"`
	From    string  `json:"from"`
	To      string  `json:"to"`
	CC      string  `json:"cc"`  // separated by comma
	BCC     string  `json:"bcc"` // separated by comma
	Subject string  `json:"subject"`
	Time    float64 `json:"time"` // lua: os.time()
	Unread  bool    `json:"unread"`
}

type Contact struct {
	Name string `json:"name"`
	Note string `json:"note"`
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
