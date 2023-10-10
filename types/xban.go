package types

type XBanRecord struct {
	Time   JsonInt `json:"time"`
	Reason string  `json:"reason"`
	Source string  `json:"source"`
}

type XBanEntry struct {
	Names   map[string]bool `json:"names"`
	Records []*XBanRecord   `json:"record"`
	Time    JsonInt         `json:"time"`
	Reason  string          `json:"reason"`
	Banned  bool            `json:"banned"`
}

// on-disk format
type XBanDatabase struct {
	Entries   []*XBanEntry `json:"entries"`
	Timestamp JsonInt      `json:"timestamp"`
}

type XBanRequest struct {
	Playername string  `json:"playername"`
	Time       JsonInt `json:"time"` //time in seconds
	Reason     string  `json:"reason"`
}
