package types

type Log struct {
	ID         string      `json:"id"`
	Timestamp  int64       `json:"timestamp"`
	Category   LogCategory `json:"category"`
	Event      LogEvent    `json:"event"`
	Username   string      `json:"username"`
	Message    string      `json:"message"`
	IPAddress  string      `json:"ip_address"`
	GeoCountry string      `json:"geo_country"`
	GeoCity    string      `json:"geo_city"`
	PosX       int         `json:"posx"`
	PosY       int         `json:"posy"`
	PosZ       int         `json:"posz"`
	Attachment []byte      `json:"attachment"`
}

type LogCategory string

const (
	CategoryUI       LogCategory = "ui"
	CategoryMinetest LogCategory = "minetest"
)

type LogEvent string

const (
	LogEventJoin  LogEvent = "join"
	LogEventLeave LogEvent = "leave"
)

type LogSearch struct {
	ID            *string      `json:"id"`
	FromTimestamp *int64       `json:"from_timestamp"`
	ToTimestamp   *int64       `json:"to_timestamp"`
	Category      *LogCategory `json:"category"`
	Event         *LogEvent    `json:"event"`
	Username      *string      `json:"username"`
	IPAddress     *string      `json:"ip_address"`
	GeoCountry    *string      `json:"geo_country"`
	Limit         *int         `json:"limit"`
}

func (m *Log) Columns(action string) []string {
	return []string{"id", "timestamp", "category", "event", "username", "message", "ip_address", "geo_country", "geo_city", "posx", "posy", "posz", "attachment"}
}

func (m *Log) Table() string {
	return "log"
}

func (m *Log) Scan(action string, r func(dest ...any) error) error {
	return r(&m.ID, &m.Timestamp, &m.Category, &m.Event, &m.Username, &m.Message, &m.IPAddress, &m.GeoCountry, &m.GeoCity, &m.PosX, &m.PosY, &m.PosZ, &m.Attachment)
}

func (m *Log) Values(action string) []any {
	return []any{m.ID, m.Timestamp, m.Category, m.Event, m.Username, m.Message, m.IPAddress, m.GeoCountry, m.GeoCity, m.PosX, m.PosY, m.PosZ, m.Attachment}
}
