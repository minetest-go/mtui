package bridge

import "encoding/json"

type CommandRequestType string
type CommandResponseType string

// UI -> Game
type CommandRequest struct {
	Type CommandRequestType `json:"type"`
	ID   *float64           `json:"id"`
	Data json.RawMessage    `json:"data"`
}

// Game -> UI
type CommandResponse struct {
	Type CommandResponseType `json:"type"`
	ID   *float64            `json:"id"`
	Data json.RawMessage     `json:"data"`
}
