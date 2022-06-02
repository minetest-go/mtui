package bridge

import "encoding/json"

type CommandType string

type Command struct {
	Type CommandType     `json:"type"`
	ID   *float64        `json:"id"`
	Data json.RawMessage `json:"data"`
}
