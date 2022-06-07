package db_test

import (
	"database/sql"
	"mtui/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*sql.DB, *db.Repositories) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	repos := db.NewRepositories(db_)
	assert.NotNil(t, db_)
	return db_, repos
}

func TestMigrate(t *testing.T) {
	setupDB(t)
}
