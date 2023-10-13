package types

type Mesecons struct {
	PosKey       string `json:"poskey"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
	Z            int    `json:"z"`
	Name         string `json:"name"`
	OrderID      int    `json:"order_id"`
	Category     string `json:"category"`
	NodeName     string `json:"nodename"`
	PlayerName   string `json:"playername"`
	State        string `json:"state"`
	LastModified int64  `json:"last_modified"`
}

func (m *Mesecons) Columns(action string) []string {
	return []string{
		"poskey",
		"x",
		"y",
		"z",
		"name",
		"order_id",
		"category",
		"nodename",
		"playername",
		"state",
		"last_modified",
	}
}

func (m *Mesecons) Table() string {
	return "mesecons"
}

func (m *Mesecons) Scan(action string, r func(dest ...any) error) error {
	return r(
		&m.PosKey,
		&m.X,
		&m.Y,
		&m.Z,
		&m.Name,
		&m.OrderID,
		&m.Category,
		&m.NodeName,
		&m.PlayerName,
		&m.State,
		&m.LastModified,
	)
}

func (m *Mesecons) Values(action string) []any {
	return []any{
		m.PosKey,
		m.X,
		m.Y,
		m.Z,
		m.Name,
		m.OrderID,
		m.Category,
		m.NodeName,
		m.PlayerName,
		m.State,
		m.LastModified,
	}
}
