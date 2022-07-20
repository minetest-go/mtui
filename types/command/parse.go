package command

import (
	"encoding/json"
	"mtui/bridge"
)

// Parse an incoming command
func ParseCommand(cmd *bridge.CommandResponse) (interface{}, error) {
	var err error
	var payload interface{}

	switch cmd.Type {
	case COMMAND_TAN_SET:
		payload = &TanCommand{}
	case COMMAND_TAN_REMOVE:
		payload = &TanCommand{}
	}

	if payload != nil {
		err = json.Unmarshal(cmd.Data, payload)
	}

	return payload, err
}
