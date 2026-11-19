// Package testutil contains shared helpers for tests
package testutil

import (
	"mtui/app"
	"mtui/types"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewApp creates an app backed by sqlite in a temporary world directory
func NewApp(t *testing.T) *app.App {
	tmpdir := t.TempDir()

	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
	`
	err := os.WriteFile(path.Join(tmpdir, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	a, err := app.Create(types.NewConfig(tmpdir))
	assert.NoError(t, err)
	assert.NotNil(t, a)
	return a
}
