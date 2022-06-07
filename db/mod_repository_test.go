package db_test

import (
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModRepository(t *testing.T) {
	_, repos := setupDB(t)

	// empty repo
	list, err := repos.Mod.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

	// create entry
	mod := &types.Mod{
		Name:       "mymod",
		ModType:    types.ModTypeRegular,
		SourceType: types.SourceTypeGit,
		URL:        "xyz",
		Version:    "123",
		AutoUpdate: true,
	}
	assert.NoError(t, repos.Mod.Save(mod))
	assert.NotEqual(t, "", mod.ID)

	// query list
	list, err = repos.Mod.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, mod.ID, list[0].ID)
	assert.Equal(t, mod.ModType, list[0].ModType)
	assert.Equal(t, mod.SourceType, list[0].SourceType)
	assert.Equal(t, mod.Name, list[0].Name)
	assert.Equal(t, mod.AutoUpdate, list[0].AutoUpdate)
	assert.Equal(t, mod.URL, list[0].URL)
	assert.Equal(t, mod.Version, list[0].Version)

	// delete
	assert.NoError(t, repos.Mod.Delete(mod.ID))
}
