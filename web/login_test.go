package web_test

import (
	"bytes"
	"encoding/json"
	"mtui/bridge"
	"mtui/types"
	"mtui/web"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenOK(t *testing.T) {
	// setup
	api, _ := CreateTestApi(t)

	// POST login

	req := &web.LoginRequest{
		Username: "singleplayer",
		Password: "mypass",
	}
	buf, err := json.Marshal(req)
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()

	api.DoLogin(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	claims := &types.Claims{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), claims))

	assert.NotNil(t, claims)
	assert.Equal(t, "singleplayer", claims.Username)
	assert.NotNil(t, claims.Privileges)
	assert.Equal(t, 1, len(claims.Privileges))
	assert.Equal(t, "interact", claims.Privileges[0])

	// GET login

	r = httptest.NewRequest("POST", "http://", nil)
	r.Header.Add("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()

	api.GetLogin(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	claims = &types.Claims{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), claims))

	assert.NotNil(t, claims)
	assert.Equal(t, "singleplayer", claims.Username)
	assert.NotNil(t, claims.Privileges)
	assert.Equal(t, 1, len(claims.Privileges))
	assert.Equal(t, "interact", claims.Privileges[0])

	// DELETE login

	r = httptest.NewRequest("DELETE", "http://", nil)
	w = httptest.NewRecorder()

	api.DoLogout(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	// Get invalid login

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	api.GetLogin(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)
}

func TestTokenTanOK(t *testing.T) {
	// setup
	api, app := CreateTestApi(t)

	// tan set
	tanset := &types.TanCommand{
		Playername: "singleplayer",
		TAN:        "1234",
	}
	tanpayload, err := json.Marshal(tanset)
	assert.NoError(t, err)

	commands := make([]*bridge.CommandResponse, 1)
	commands[0] = &bridge.CommandResponse{
		Type: types.COMMAND_TAN_SET,
		Data: tanpayload,
	}

	buf, err := json.Marshal(commands)
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()

	app.Bridge.HandlePost(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	time.Sleep(time.Millisecond * 20)

	// POST login

	req := &web.LoginRequest{
		Username: "singleplayer",
		Password: "1234",
	}
	buf, err = json.Marshal(req)
	assert.NoError(t, err)

	r = httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()

	api.DoLogin(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	claims := &types.Claims{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), claims))
}

func TestTokenFailed(t *testing.T) {
	// setup
	api, _ := CreateTestApi(t)

	// GET login

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetLogin(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)
}
