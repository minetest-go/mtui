package web_test

import (
	"mtui/app"
	"mtui/auth"
	"mtui/testutil"
	"mtui/web"
	"testing"

	authdb "github.com/luanti-golang/mtdb/auth"
	"github.com/stretchr/testify/assert"
)

func CreateTestApi(t *testing.T) (*web.Api, *app.App) {
	app := testutil.NewApp(t)

	api := web.NewApi(app)
	assert.NoError(t, api.Setup())

	assert.NotNil(t, api)

	// create test data

	salt, verifier, err := auth.CreateAuth("singleplayer", "mypass")
	assert.NoError(t, err)

	dbpass := auth.CreateDBPassword(salt, verifier)

	// create user

	auth_entry := &authdb.AuthEntry{
		Name:      "singleplayer",
		Password:  dbpass,
		LastLogin: 123,
	}
	assert.NoError(t, app.DBContext.Auth.Create(auth_entry))
	assert.NotNil(t, auth_entry.ID)

	// create privs

	assert.NoError(t, app.DBContext.Privs.Create(&authdb.PrivilegeEntry{
		ID:        *auth_entry.ID,
		Privilege: "interact",
	}))

	return api, app
}
