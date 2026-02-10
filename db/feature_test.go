package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeatureRepo(t *testing.T) {
	_, g := setupDB(t)
	repo := db.NewRepositories(g).FeatureRepository

	f1_name := types.FeatureName("f1")
	f2_name := types.FeatureName("f2")

	// create
	assert.NoError(t, repo.Set(&types.Feature{Name: f1_name, Enabled: true}))
	assert.NoError(t, repo.Set(&types.Feature{Name: f2_name, Enabled: false}))

	// read
	list, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))

	f1, err := repo.GetByName(f1_name)
	assert.NoError(t, err)
	assert.Equal(t, f1_name, f1.Name)
	assert.Equal(t, true, f1.Enabled)

	f2, err := repo.GetByName(f2_name)
	assert.NoError(t, err)
	assert.Equal(t, f2_name, f2.Name)
	assert.Equal(t, false, f2.Enabled)

	// update
	assert.NoError(t, repo.Set(&types.Feature{Name: f1_name, Enabled: false}))

	// read
	f1, err = repo.GetByName(f1_name)
	assert.NoError(t, err)
	assert.Equal(t, f1_name, f1.Name)
	assert.Equal(t, false, f1.Enabled)
}
