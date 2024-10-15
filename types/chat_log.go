package types

type ChatLog struct {
	ID        string `json:"id" gorm:"primarykey;column:id"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
	Channel   string `json:"channel" gorm:"column:channel"`
	Name      string `json:"name" gorm:"column:name"`
	Message   string `json:"message" gorm:"column:message"`
}

func (m *ChatLog) TableName() string {
	return "chat_log"
}
