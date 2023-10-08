package types

type Log struct {
	ID         string      `json:"id"`
	Timestamp  int64       `json:"timestamp"`
	Category   LogCategory `json:"category"`
	Event      string      `json:"event"`
	Username   string      `json:"username"`
	Message    string      `json:"message"`
	Nodename   *string     `json:"nodename"`
	IPAddress  *string     `json:"ip_address"`
	GeoCountry *string     `json:"geo_country"`
	GeoCity    *string     `json:"geo_city"`
	GeoASN     *int        `json:"geo_asn"`
	PosX       *JsonInt    `json:"posx"`
	PosY       *JsonInt    `json:"posy"`
	PosZ       *JsonInt    `json:"posz"`
	Attachment []byte      `json:"attachment"`
}

type LogCategory string

const (
	CategoryUI       LogCategory = "ui"
	CategoryMinetest LogCategory = "minetest"
	CategoryService  LogCategory = "service"
)

type LogSearch struct {
	ID            *string      `json:"id"`
	FromTimestamp *int64       `json:"from_timestamp"`
	ToTimestamp   *int64       `json:"to_timestamp"`
	Category      *LogCategory `json:"category"`
	Event         *string      `json:"event"`
	Username      *string      `json:"username"`
	IPAddress     *string      `json:"ip_address"`
	GeoCountry    *string      `json:"geo_country"`
	GeoCity       *string      `json:"geo_city"`
	GeoASN        *int         `json:"geo_asn"`
	Limit         *int         `json:"limit"`
	PosX          *int         `json:"posx"`
	PosY          *int         `json:"posy"`
	PosZ          *int         `json:"posz"`
	Nodename      *string      `json:"nodename"`
	SortAscending bool         `json:"sort_ascending"`
}

func (m *Log) Columns(action string) []string {
	return []string{
		"id",
		"timestamp",
		"category",
		"event",
		"username",
		"message",
		"nodename",
		"ip_address",
		"geo_country",
		"geo_city",
		"geo_asn",
		"posx",
		"posy",
		"posz",
		"attachment",
	}
}

func (m *Log) Table() string {
	return "log"
}

func (m *Log) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.ID,
		&m.Timestamp,
		&m.Category,
		&m.Event,
		&m.Username,
		&m.Message,
		&m.Nodename,
		&m.IPAddress,
		&m.GeoCountry,
		&m.GeoCity,
		&m.GeoASN,
		&m.PosX,
		&m.PosY,
		&m.PosZ,
		&m.Attachment,
	)
}

func (m *Log) Values(action string) []any {
	return []any{
		m.ID,
		m.Timestamp,
		m.Category,
		m.Event,
		m.Username,
		m.Message,
		m.Nodename,
		m.IPAddress,
		m.GeoCountry,
		m.GeoCity,
		m.GeoASN,
		m.PosX,
		m.PosY,
		m.PosZ,
		m.Attachment,
	}
}
