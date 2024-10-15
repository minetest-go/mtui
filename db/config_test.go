package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	DB, g := setupDB(t)
	repo := db.NewRepositories(DB, g).ConfigRepo

	// create
	assert.NoError(t, repo.Set(&types.ConfigEntry{Key: "x", Value: "y"}))

	// read
	c, err := repo.GetByKey("x")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "y", c.Value)

	// update
	assert.NoError(t, repo.Set(&types.ConfigEntry{Key: "x", Value: "y2"}))

	// read
	c, err = repo.GetByKey("x")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "y2", c.Value)

	// read
	c, err = repo.GetByKey("z")
	assert.NoError(t, err)
	assert.Nil(t, c)

	// delete
	assert.NoError(t, repo.Delete("x"))

	// read
	c, err = repo.GetByKey("x")
	assert.NoError(t, err)
	assert.Nil(t, c)
}
