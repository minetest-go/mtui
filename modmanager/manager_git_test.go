package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckoutBranch(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	// checkout master
	mod := &types.Mod{
		Name:       "moreblocks",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeGIT,
		URL:        "https://github.com/minetest-mods/moreblocks.git",
		Branch:     "master",
	}
	assert.NoError(t, mm.Create(mod))
	assert.True(t, mod.Version != "")

	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))
}

func TestCheckoutHash(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	// checkout master branch on specified commit
	mod := &types.Mod{
		Name:       "moreblocks",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeGIT,
		URL:        "https://github.com/minetest-mods/moreblocks.git",
		Branch:     "master",
		Version:    "fe34e3f3cd3e066ba0be76f9df46c11e66411496",
	}
	assert.NoError(t, mm.Create(mod))
	assert.NotEqual(t, "", mod.ID)
	assert.Equal(t, "fe34e3f3cd3e066ba0be76f9df46c11e66411496", mod.Version)

	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))
	assert.Equal(t, "moreblocks", mods[0].Name)
	assert.Equal(t, "https://github.com/minetest-mods/moreblocks.git", mods[0].URL)
	assert.Equal(t, "master", mods[0].Branch)
	assert.Equal(t, "fe34e3f3cd3e066ba0be76f9df46c11e66411496", mods[0].Version)
	assert.Equal(t, types.SourceTypeGIT, mods[0].SourceType)
	assert.Equal(t, types.ModTypeMod, mods[0].ModType)
	mod = mods[0]

	// check list
	mods, err = app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 1, len(mods))

	// check remote status
	err = mm.CheckUpdates()
	assert.NoError(t, err)

	mod, err = app.Repos.ModRepo.GetByID(mod.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, "", mod.LatestVersion)
	assert.NotEqual(t, "fe34e3f3cd3e066ba0be76f9df46c11e66411496", mod.LatestVersion)

	// update
	assert.NoError(t, mm.Update(mod, mod.LatestVersion))

	// remove
	assert.NoError(t, mm.Remove(mod))
	mods, err = app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, mods)
	assert.Equal(t, 0, len(mods))

}

func TestCheckoutGame(t *testing.T) {
	//t.Skip() // slow test, enable on demand

	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	// checkout master branch on specified commit
	mod := &types.Mod{
		Name:       "game",
		ModType:    types.ModTypeGame,
		SourceType: types.SourceTypeGIT,
		URL:        "https://github.com/ubc-minetest-classroom/minetest_classroom",
		Branch:     "main",
	}
	assert.NoError(t, mm.Create(mod))
	assert.NotEqual(t, "", mod.ID)
	assert.True(t, mod.Version != "")

	// remove
	assert.NoError(t, mm.Remove(mod))
}
