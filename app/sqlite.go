package app

import (
	"context"
	"fmt"
	"mtui/db"
	"os"
	"strings"
)

func IgnoreSqliteFileDownload(filename string) bool {
	if strings.HasSuffix(filename, ".sqlite-shm") || strings.HasSuffix(filename, ".sqlite-wal") {
		// sqlite wal and shared memory
		return true
	}
	return false
}

func IsSqliteDatabase(filename string) bool {
	return strings.HasSuffix(filename, ".sqlite")
}

func CreateSqliteSnapshot(filename string) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "backup.sqlite")
	if err != nil {
		return "", fmt.Errorf("create temp error: %v", err)
	}

	err = db.BackupSqlite3Database(context.Background(), filename, f.Name())
	if err != nil {
		return "", fmt.Errorf("backup error from '%s' to '%s': %v", filename, f.Name(), err)
	}

	return f.Name(), nil
}
