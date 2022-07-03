package modmanager_test

import (
	"mtui/modmanager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)
	assert.NotNil(t, mm)

	err := mm.Scan()
	assert.NoError(t, err)

	mods := mm.Mods()
	assert.NotNil(t, mods)
	assert.Equal(t, 0, len(mods))

	mod := &modmanager.Mod{
		Name:       "moreblocks",
		ModType:    modmanager.ModTypeRegular,
		SourceType: modmanager.SourceTypeGit,
		URL:        "https://github.com/minetest-mods/moreblocks.git",
		Version:    "fe34e3f3cd3e066ba0be76f9df46c11e66411496",
	}

	assert.NoError(t, mm.Create(mod))

	status, err := mm.Status(mod)
	assert.NoError(t, err)
	assert.NotNil(t, status)

	assert.NoError(t, mm.Remove(mod))

	// TODO
}
