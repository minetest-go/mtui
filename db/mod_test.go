package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestModRepo(t *testing.T) {
	_, g := setupDB(t)
	repo := db.NewRepositories(g).ModRepo

	m := &types.Mod{
		ID:         uuid.NewString(),
		Name:       "my_mod",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeGIT,
		URL:        "https://",
		Version:    "1234",
		AutoUpdate: false,
	}

	// Create
	assert.NoError(t, repo.Create(m))

	// Read
	list, err := repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, m.ID, list[0].ID)
	assert.Equal(t, m.ModType, list[0].ModType)
	assert.Equal(t, m.Name, list[0].Name)
	assert.Equal(t, m.URL, list[0].URL)

	// by id
	m2, err := repo.GetByID(m.ID)
	assert.NoError(t, err)
	assert.NotNil(t, m2)

	// by non-existent id
	m2, err = repo.GetByID(uuid.NewString())
	assert.NoError(t, err)
	assert.Nil(t, m2)

	// Update
	m.URL = "xyz"
	assert.NoError(t, repo.Update(m))

	list, err = repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, m.URL, list[0].URL)

	// Delete
	assert.NoError(t, repo.Delete(m.ID))

	list, err = repo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

}
