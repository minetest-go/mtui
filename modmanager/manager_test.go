package modmanager_test

import (
	"mtui/modmanager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)
	assert.NotNil(t, mm)

	err := mm.Scan()
	assert.NoError(t, err)

	mods := mm.Mods()
	assert.NotNil(t, mods)
	assert.Equal(t, 0, len(mods))
}

func TestCheckoutBranch(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)

	// checkout master
	mod := &modmanager.Mod{
		Name:       "moreblocks",
		ModType:    modmanager.ModTypeRegular,
		SourceType: modmanager.SourceTypeGit,
		URL:        "https://github.com/minetest-mods/moreblocks.git",
		Branch:     "refs/heads/master",
	}
	assert.NoError(t, mm.Create(mod))
	assert.True(t, mod.Version != "")

	mods := mm.Mods()
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))
}

func TestCheckoutHash(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)

	// checkout master
	mod := &modmanager.Mod{
		Name:       "moreblocks",
		ModType:    modmanager.ModTypeRegular,
		SourceType: modmanager.SourceTypeGit,
		URL:        "https://github.com/minetest-mods/moreblocks.git",
		Branch:     "refs/heads/master",
		Version:    "fe34e3f3cd3e066ba0be76f9df46c11e66411496",
	}
	assert.NoError(t, mm.Create(mod))
	assert.Equal(t, "fe34e3f3cd3e066ba0be76f9df46c11e66411496", mod.Version)

	mods := mm.Mods()
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))

	status, err := mm.Status(mod)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "fe34e3f3cd3e066ba0be76f9df46c11e66411496", status.CurrentVersion)
	assert.True(t, status.LatestVersion != "")
	assert.True(t, status.LatestVersion != status.CurrentVersion)

	assert.NoError(t, mm.Update(mod, status.LatestVersion))
	status2, err := mm.Status(mod)
	assert.NoError(t, err)
	assert.NotNil(t, status2)
	assert.Equal(t, status.LatestVersion, status2.CurrentVersion)
	assert.Equal(t, status.LatestVersion, status2.LatestVersion)
}
