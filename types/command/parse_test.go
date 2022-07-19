package command_test

import (
	"mtui/bridge"
	"mtui/types/command"
	"testing"

	"github.com/stretchr/testify/assert"
)

const UNKNOWN_RES bridge.CommandResponseType = "whatever"

func TestParseCommand(t *testing.T) {
	resp := &bridge.CommandResponse{
		Type: command.COMMAND_TAN_REMOVE,
		Data: []byte("{}"),
	}

	o, err := command.ParseCommand(resp)
	assert.NoError(t, err)
	assert.NotNil(t, o)

	resp = &bridge.CommandResponse{
		Type: UNKNOWN_RES,
		Data: []byte("{}"),
	}

	o, err = command.ParseCommand(resp)
	assert.NoError(t, err)
	assert.Nil(t, o)
}
