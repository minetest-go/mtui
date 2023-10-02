package types

type ChatLog struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Channel   string `json:"channel"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}

func (m *ChatLog) Columns(action string) []string {
	return []string{
		"id",
		"timestamp",
		"channel",
		"name",
		"message",
	}
}

func (m *ChatLog) Table() string {
	return "chat_log"
}

func (m *ChatLog) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Timestamp,
		&m.Channel,
		&m.Name,
		&m.Message,
	)
}

func (m *ChatLog) Values(action string) []any {
	return []any{
		m.ID,
		m.Timestamp,
		m.Channel,
		m.Name,
		m.Message,
	}
}
