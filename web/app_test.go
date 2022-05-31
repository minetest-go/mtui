package web_test

import (
	"mtadmin/app"
	"mtadmin/web"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestApp(t *testing.T) *app.App {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtadmin_app")
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
	os.Setenv("JWTKEY", "mykey")
	app := CreateTestApp(t)

	api := web.NewApi(app)
	assert.NotNil(t, api)

	return api, app
}
