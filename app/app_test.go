package app_test

import (
	"mtadmin/app"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateApp(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtadmin_app")
	assert.NoError(t, err)

	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
	`
	err = os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	app, err := app.Create(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}
