package mail

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
	CC      *string `json:"cc"`  // separated by comma
	BCC     *string `json:"bcc"` // separated by comma
	Subject string  `json:"subject"`
	Time    float64 `json:"time"` // lua: os.time()
	Read    bool    `json:"read"`
}

type Contact struct {
	Name string `json:"name"`
	Note string `json:"note"`
}
