package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapRepo(t *testing.T) {
	// temp files
	dbfile, err := os.CreateTemp(os.TempDir(), "map.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)

	// create repo
	repo, err := NewMapRepository(dbfile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	// size
	size, err := repo.GetSize()
	assert.NoError(t, err)
	assert.True(t, size > 0)

	// count
	count, err := repo.CountBlocks()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
