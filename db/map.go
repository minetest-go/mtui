package db

import "database/sql"

const MAP_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists blocks (
	pos INT PRIMARY KEY,
	data BLOB
);
`

func NewMapRepository(filename string) (*MapRepository, error) {
	db_, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	_, err = db_.Exec(MAP_MIGRATE)
	if err != nil {
		return nil, err
	}

	return &MapRepository{db: db_}, nil
}

type MapRepository struct {
	db *sql.DB
}
