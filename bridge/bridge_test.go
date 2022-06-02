package bridge_test

import (
	"bytes"
	"encoding/json"
	"mtui/bridge"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestCommand bridge.CommandType = "test"

func TestBridgeSendCommand(t *testing.T) {
	b := bridge.New()
	b.SendCommand(TestCommand, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	b.HandleGet(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	var cmds []*bridge.Command
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&cmds))
	assert.NotNil(t, cmds)
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, TestCommand, cmds[0].Type)
}

func TestBridgeReceiveCommand(t *testing.T) {
	b := bridge.New()
	c := b.AddHandler()

	var cmd *bridge.Command
	select {
	case cmd = <-c:
	default:
	}
	assert.Nil(t, cmd)

	buf, err := json.Marshal(&bridge.Command{Type: TestCommand})
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()

	b.HandlePost(w, r)

	cmd = <-c

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.NotNil(t, cmd)
	assert.Equal(t, TestCommand, cmd.Type)
}
