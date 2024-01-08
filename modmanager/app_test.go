package modmanager_test

import (
	"mtui/app"
	"mtui/types"
	"os"
	"path"
	"testing"

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

	a, err := app.Create(types.NewConfig(tmpdir))
	assert.NoError(t, err)
	return a
}
