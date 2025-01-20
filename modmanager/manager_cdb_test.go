package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestCDBRelease(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)

	// checkout master
	mod := &types.Mod{
		Name:       "blockexchange",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeCDB,
		Author:     "BuckarooBanzay",
	}
	assert.NoError(t, mm.Create(mod))
	assert.NoError(t, app.Repos.ModRepo.Create(mod))
	assert.True(t, mod.Version != "")

	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))

	changed_mods, err := modmanager.CheckUpdates(app.WorldDir, mods)
	assert.NoError(t, err)
	assert.NotNil(t, changed_mods)

	mod, err = app.Repos.ModRepo.GetByID(mod.ID)
	assert.NoError(t, err)
	assert.Equal(t, mod.Version, mod.LatestVersion)

	err = mm.Update(mod, mod.LatestVersion)
	assert.NoError(t, err)

	err = mm.Remove(mod)
	assert.NoError(t, err)
}
