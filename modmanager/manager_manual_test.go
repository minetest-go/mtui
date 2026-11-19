package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManualMod(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir, app.Repos.ModRepo)

	dir := path.Join(app.WorldDir, "worldmods", "legacy")
	assert.NoError(t, os.MkdirAll(dir, 0755))

	mod := &types.Mod{Name: "legacy", ModType: types.ModTypeMod, SourceType: types.SourceTypeManual, Status: types.ModStatusInstalled}
	assert.NoError(t, app.Repos.ModRepo.Create(mod))

	// no updates, no panic
	assert.NoError(t, mm.CheckUpdates())
	assert.Error(t, mm.Update(mod, "1"))

	// removal deletes the directory and the entry
	assert.NoError(t, mm.Remove(mod))
	_, err := os.Stat(dir)
	assert.True(t, os.IsNotExist(err))
	mods, err := app.Repos.ModRepo.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, mods)

	// unknown source types are rejected instead of panicking
	assert.Error(t, mm.Create(&types.Mod{Name: "x", SourceType: "cdb"}))
}
