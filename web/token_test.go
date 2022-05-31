package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mtadmin/auth"
	"mtadmin/types"
	"mtadmin/web"
	"net/http/httptest"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestTokenOK(t *testing.T) {
	// setup
	api, app := CreateTestApi(t)

	salt, verifier, err := auth.CreateAuth("singleplayer", "mypass")
	assert.NoError(t, err)

	dbpass := auth.CreateDBPassword(salt, verifier)

	// create user

	auth_entry := &mtdb.AuthEntry{
		Name:      "singleplayer",
		Password:  dbpass,
		LastLogin: 123,
	}
	assert.NoError(t, app.DBContext.Auth.Create(auth_entry))
	assert.NotNil(t, auth_entry.ID)

	// create privs

	assert.NoError(t, app.DBContext.Privs.Create(&mtdb.PrivilegeEntry{
		ID:        *auth_entry.ID,
		Privilege: "interact",
	}))

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

	// GET login

	fmt.Printf("%s\n", w.Header().Get("Set-Cookie"))
	r = httptest.NewRequest("POST", "http://", nil)
	r.Header.Add("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()

	api.GetLogin(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	claims := &types.Claims{}
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

func TestTokenFailed(t *testing.T) {
	// setup
	api, _ := CreateTestApi(t)

	// GET login

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetLogin(w, r)

	assert.Equal(t, 401, w.Result().StatusCode)
}
