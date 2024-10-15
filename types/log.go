package types

type Log struct {
	ID         string      `json:"id" gorm:"primarykey;column:id"`
	Timestamp  int64       `json:"timestamp" gorm:"column:timestamp"`
	Category   LogCategory `json:"category" gorm:"column:category"`
	Event      string      `json:"event" gorm:"column:event"`
	Username   string      `json:"username" gorm:"column:username"`
	Message    string      `json:"message" gorm:"column:message"`
	Nodename   *string     `json:"nodename" gorm:"column:nodename"`
	IPAddress  *string     `json:"ip_address" gorm:"column:ip_address"`
	GeoCountry *string     `json:"geo_country" gorm:"column:geo_country"`
	GeoCity    *string     `json:"geo_city" gorm:"column:geo_city"`
	GeoASN     *int        `json:"geo_asn" gorm:"column:geo_asn"`
	PosX       *JsonInt    `json:"posx" gorm:"column:posx"`
	PosY       *JsonInt    `json:"posy" gorm:"column:posy"`
	PosZ       *JsonInt    `json:"posz" gorm:"column:posz"`
	Attachment []byte      `json:"attachment" gorm:"column:attachment"`
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

func (m *Log) TableName() string {
	return "log"
}
