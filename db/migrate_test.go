package db_test

import (
	"database/sql"
	"mtui/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *sql.DB {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db_
}

func TestMigrate(t *testing.T) {
	setupDB(t)
}
