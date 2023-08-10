package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
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
		Author:     "buckaroobanzay",
	}
	assert.NoError(t, mm.Create(mod))
	assert.True(t, mod.Version != "")

	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))

	status, err := mm.Status(mod)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, mod.Version, status.CurrentVersion)
	assert.Equal(t, mod.Version, status.LatestVersion)

	err = mm.Update(mod, mod.Version)
	assert.NoError(t, err)

	err = mm.Remove(mod)
	assert.NoError(t, err)
}
