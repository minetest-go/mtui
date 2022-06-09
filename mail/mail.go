package mail

import (
	"encoding/json"
	"errors"
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

func (m *Mail) getMailFile(playername string) string {
	return path.Join(m.world_dir, "mails", playername+".json")
}

func (m *Mail) getContactsFile(playername string) string {
	return path.Join(m.world_dir, "mails", "contacts", playername+".json")
}

func (m *Mail) GetMessages(playername string) ([]*Message, error) {
	b, err := os.ReadFile(m.getMailFile(playername))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	list := make([]*Message, 0)
	err = json.Unmarshal(b, &list)
	return list, err
}

func (m *Mail) SetMessages(playername string, list []*Message) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(m.getMailFile(playername), b, 0755)
}

func (m *Mail) GetContacts(playername string) (map[string]*Contact, error) {
	b, err := os.ReadFile(m.getContactsFile(playername))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	c := make(map[string]*Contact)
	err = json.Unmarshal(b, &c)
	return c, err
}
