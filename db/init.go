package db

import (
	"database/sql"
	"fmt"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(world_dir string) (*sql.DB, *gorm.DB, error) {

	connStr := "mtui.sqlite?_timeout=5000&_journal=WAL&_cache=shared"
	var err error
	db, err := sql.Open("sqlite3", path.Join(world_dir, connStr))
	if err != nil {
		return nil, nil, fmt.Errorf("open error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("ping error: %v", err)
	}

	g, err := gorm.Open(sqlite.Open(path.Join(world_dir, connStr)))
	if err != nil {
		return nil, nil, fmt.Errorf("gorm.Open error: %v", err)
	}

	return db, g, nil
}
