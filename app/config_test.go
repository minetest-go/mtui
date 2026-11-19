package app_test

import (
	"mtui/testutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMTConfigCreatesReadableFile(t *testing.T) {
	a := testutil.NewApp(t)
	a.Config.LuantiConfig = path.Join(a.WorldDir, "minetest.conf")

	s, err := a.ReadMTConfig(nil)
	assert.NoError(t, err)
	assert.Empty(t, s)

	fi, err := os.Stat(a.Config.LuantiConfig)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), fi.Mode().Perm())
}

func TestReadMTConfigRepairsPermissions(t *testing.T) {
	a := testutil.NewApp(t)
	a.Config.LuantiConfig = path.Join(a.WorldDir, "minetest.conf")

	// file as created by earlier versions
	assert.NoError(t, os.WriteFile(a.Config.LuantiConfig, []byte("max_users = 10\n"), 0644))
	assert.NoError(t, os.Chmod(a.Config.LuantiConfig, 0))

	s, err := a.ReadMTConfig(nil)
	assert.NoError(t, err)
	assert.Equal(t, "10", s["max_users"].Value)

	fi, err := os.Stat(a.Config.LuantiConfig)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), fi.Mode().Perm())
}
