package db_test

import (
	"database/sql"
	"mtui/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*sql.DB, *gorm.DB) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	db_, g, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)
	assert.NotNil(t, g)

	err = db.Migrate(db_)
	assert.NoError(t, err)

	return db_, g
}

func TestMigrate(t *testing.T) {
	setupDB(t)
}
