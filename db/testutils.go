package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*sql.DB, *Repositories) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	db_, err := Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = Migrate(db_)
	assert.NoError(t, err)

	repos := NewRepositories(db_)
	assert.NotNil(t, db_)
	return db_, repos
}
