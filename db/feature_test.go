package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeatureRepo(t *testing.T) {
	_db, g := setupDB(t)
	repo := db.NewRepositories(_db, g).FeatureRepository
	// create
	assert.NoError(t, repo.Set(&types.Feature{Name: "f1", Enabled: true}))
	assert.NoError(t, repo.Set(&types.Feature{Name: "f2", Enabled: false}))

	// read
	list, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))

	f1, err := repo.GetByName("f1")
	assert.NoError(t, err)
	assert.Equal(t, "f1", f1.Name)
	assert.Equal(t, true, f1.Enabled)

	f2, err := repo.GetByName("f2")
	assert.NoError(t, err)
	assert.Equal(t, "f2", f2.Name)
	assert.Equal(t, false, f2.Enabled)

	// update
	assert.NoError(t, repo.Set(&types.Feature{Name: "f1", Enabled: false}))

	// read
	f1, err = repo.GetByName("f1")
	assert.NoError(t, err)
	assert.Equal(t, "f1", f1.Name)
	assert.Equal(t, false, f1.Enabled)
}
