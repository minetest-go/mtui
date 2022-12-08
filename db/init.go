package db

import (
	"database/sql"
	"path"

	"github.com/minetest-go/mtdb/wal"
	_ "modernc.org/sqlite"
)

func Init(world_dir string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite", path.Join(world_dir, "mtui.sqlite"))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = wal.EnableWAL(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
