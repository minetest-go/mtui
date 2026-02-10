package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestCDBRelease(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	// checkout master
	mod := &types.Mod{
		Name:       "blockexchange",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeCDB,
		Author:     "BuckarooBanzay",
	}
	assert.NoError(t, mm.Create(mod))
	assert.True(t, mod.Version != "")

	fi, err := os.Stat(path.Join(app.WorldDir, "worldmods", "blockexchange"))
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.True(t, fi.IsDir())

	fi, err = os.Stat(path.Join(app.WorldDir, "worldmods", "blockexchange", "init.lua"))
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.False(t, fi.IsDir())

	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))

	err = mm.CheckUpdates()
	assert.NoError(t, err)

	mod, err = app.Repos.ModRepo.GetByID(mod.ID)
	assert.NoError(t, err)
	assert.Equal(t, mod.Version, mod.LatestVersion)

	err = mm.Update(mod, mod.LatestVersion)
	assert.NoError(t, err)

	err = mm.Remove(mod)
	assert.NoError(t, err)
}

func TestShortnamesMod(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	// checkout master
	mod := &types.Mod{
		Name:       "shortnames",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeCDB,
		Author:     "niwla23",
	}
	assert.NoError(t, mm.Create(mod))
	assert.True(t, mod.Version != "")

	fi, err := os.Stat(path.Join(app.WorldDir, "worldmods", "shortnames"))
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.True(t, fi.IsDir())
}
