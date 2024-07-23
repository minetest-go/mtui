package db_test

import (
	"context"
	"mtui/db"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackupSqlite3(t *testing.T) {
	DB := setupDB(t)

	dstfile, err := os.CreateTemp(os.TempDir(), "backup-dest.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dstfile)

	err = db.BackupSqlite3Database(context.Background(), DB, dstfile.Name())
	assert.NoError(t, err)
}
