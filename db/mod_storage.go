package db

import "database/sql"

const MOD_STORAGE_MIGRATE = `
pragma journal_mode = wal;
begin;
CREATE TABLE if not exists entries (
	modname TEXT NOT NULL,
	key BLOB NOT NULL,
	value BLOB NOT NULL,
	PRIMARY KEY (modname, key)
);
commit;
`

func NewModStorageRepository(filename string) (*ModStorageRepository, error) {
	db_, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	_, err = db_.Exec(MOD_STORAGE_MIGRATE)
	if err != nil {
		return nil, err
	}

	return &ModStorageRepository{db: db_}, nil
}

type ModStorageRepository struct {
	db *sql.DB
}
