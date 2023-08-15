package db

import (
	"database/sql"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func Init(world_dir string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite3", path.Join(world_dir, "mtui.sqlite?_timeout=5000&_journal=WAL&_cache=shared"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
