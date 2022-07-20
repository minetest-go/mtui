package bridge

import "encoding/json"

type CommandType string

// UI -> Game
type CommandRequest struct {
	Type CommandType     `json:"type"`
	ID   *float64        `json:"id"`
	Data json.RawMessage `json:"data"`
}

// Game -> UI
type CommandResponse struct {
	Type CommandType     `json:"type"`
	ID   *float64        `json:"id"`
	Data json.RawMessage `json:"data"`
}
