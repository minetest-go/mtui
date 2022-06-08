package mail

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
	//TODO
	return nil, nil
}

func (m *Mail) Contacts(playername string) (map[string]*Contact, error) {
	//TODO
	return nil, nil
}
