package bridge_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TestCommandRequest  bridge.CommandRequestType  = "test"
	TestCommandResponse bridge.CommandResponseType = "test"
)

func TestBridgeSendCommand(t *testing.T) {
	b := bridge.New()
	b.SendCommand(TestCommandRequest, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	b.HandleGet(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	var cmds []*bridge.CommandRequest
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&cmds))
	assert.NotNil(t, cmds)
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, TestCommandRequest, cmds[0].Type)
}

func TestBridgeReceiveCommand(t *testing.T) {
	b := bridge.New()
	c := b.AddHandler()

	var cmd *bridge.CommandResponse
	select {
	case cmd = <-c:
	default:
	}
	assert.Nil(t, cmd)

	commands := make([]*bridge.CommandRequest, 1)
	commands[0] = &bridge.CommandRequest{Type: TestCommandRequest}
	buf, err := json.Marshal(commands)
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()

	b.HandlePost(w, r)

	cmd = <-c

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.NotNil(t, cmd)
	assert.Equal(t, TestCommandResponse, cmd.Type)
}

func TestBridgeReceiveInvalidCommand(t *testing.T) {
	b := bridge.New()

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer([]byte("blah")))
	w := httptest.NewRecorder()

	b.HandlePost(w, r)
	assert.Equal(t, 500, w.Result().StatusCode)
}

func TestBridgeExecuteCommandTimeout(t *testing.T) {
	b := bridge.New()
	cmd, err := b.ExecuteCommand(TestCommandRequest, nil, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Nil(t, cmd)
}

func TestBridgeExecuteCommand(t *testing.T) {
	b := bridge.New()
	var rx_cmd *bridge.CommandResponse
	var rx_err error

	go func() {
		rx_cmd, rx_err = b.ExecuteCommand(TestCommandRequest, nil, 500*time.Millisecond)
		fmt.Printf("execution finished: id=%f\n", *rx_cmd.ID)
	}()

	// get command from bridge
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	b.HandleGet(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	var cmds []*bridge.CommandRequest
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&cmds))
	assert.NotNil(t, cmds)
	assert.Equal(t, 1, len(cmds))
	assert.Equal(t, TestCommandRequest, cmds[0].Type)
	assert.NotNil(t, cmds[0].ID)
	assert.True(t, *cmds[0].ID > 0)

	// send response to bridge
	commands := make([]*bridge.CommandResponse, 1)
	commands[0] = &bridge.CommandResponse{Type: TestCommandResponse, ID: cmds[0].ID}
	buf, err := json.Marshal(commands)
	assert.NoError(t, err)

	r = httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()

	b.HandlePost(w, r)

	time.Sleep(200 * time.Millisecond)

	// assert result
	assert.Nil(t, rx_err)
	assert.NotNil(t, rx_cmd)
}
