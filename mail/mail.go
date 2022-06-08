package mail

import (
	"encoding/json"
	"os"
	"path"
)

type Mail struct {
	world_dir string
}

func New(world_dir string) *Mail {
	return &Mail{world_dir: world_dir}
}

type Message struct {
	Body    string `json:"body"`
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Time    int64  `json:"time"`
	Unread  bool   `json:"unread"`
}

type Contact struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func (m *Mail) List(playername string) ([]*Message, error) {
	f, err := os.Open(path.Join(m.world_dir, "mails", playername+".json"))
	if err != nil {
		return nil, err
	}

	list := make([]*Message, 0)
	err = json.NewDecoder(f).Decode(&list)
	return list, nil
}

func (m *Mail) Contacts(playername string) (map[string]*Contact, error) {
	f, err := os.Open(path.Join(m.world_dir, "mails", "contacts", playername+".json"))
	if err != nil {
		return nil, err
	}

	c := make(map[string]*Contact)
	err = json.NewDecoder(f).Decode(&c)
	return c, nil
}
