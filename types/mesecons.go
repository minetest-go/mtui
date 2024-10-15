package types

type Mesecons struct {
	PosKey       string `json:"poskey" gorm:"primarykey;column:poskey"`
	X            int    `json:"x" gorm:"column:x"`
	Y            int    `json:"y" gorm:"column:y"`
	Z            int    `json:"z" gorm:"column:z"`
	Name         string `json:"name" gorm:"column:name"`
	OrderID      int    `json:"order_id" gorm:"column:order_id"`
	Category     string `json:"category" gorm:"column:category"`
	NodeName     string `json:"nodename" gorm:"column:nodename"`
	PlayerName   string `json:"playername" gorm:"column:playername"`
	State        string `json:"state" gorm:"column:state"`
	LastModified int64  `json:"last_modified" gorm:"column:last_modified"`
}

func (m *Mesecons) TableName() string {
	return "mesecons"
}
