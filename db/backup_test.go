package db_test

import (
	"context"
	"mtui/db"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackupSqlite3(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	db_, err := db.Init(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, db_)

	err = db.Migrate(db_)
	assert.NoError(t, err)
	dstfile, err := os.CreateTemp(os.TempDir(), "backup-dest.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dstfile)

	err = db.BackupSqlite3Database(context.Background(), path.Join(tmpdir, "mtui.sqlite"), dstfile.Name())
	assert.NoError(t, err)
}
