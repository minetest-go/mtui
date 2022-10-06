package web_test

import (
	"mtui/app"
	"mtui/auth"
	"mtui/web"
	"os"
	"path"
	"testing"

	authdb "github.com/minetest-go/mtdb/auth"
	"github.com/stretchr/testify/assert"
)

func CreateTestApp(t *testing.T) *app.App {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui_app")
	assert.NoError(t, err)

	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
	`
	err = os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	a, err := app.Create(tmpdir)
	assert.NoError(t, err)
	return a
}

func CreateTestApi(t *testing.T) (*web.Api, *app.App) {
	app := CreateTestApp(t)

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
