package command

import "mtui/bridge"

type SentMessage struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	CC      *string `json:"cc"`  // separated by comma
	BCC     *string `json:"bcc"` // separated by comma
	Subject string  `json:"subject"`
	Body    string  `json:"body"`
}

const (
	COMMAND_MAIL_SENT bridge.CommandType = "mail_sent"
)
