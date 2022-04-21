package db

import "database/sql"

const AUTH_MIGRATE = `
pragma journal_mode = wal;
begin;
CREATE TABLE if not exists auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
CREATE TABLE if not exists user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
commit;
`

func NewAuthRepository(filename string) (*AuthRepository, error) {
	db_, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	_, err = db_.Exec(AUTH_MIGRATE)
	if err != nil {
		return nil, err
	}

	return &AuthRepository{db: db_}, nil
}

type AuthRepository struct {
	db *sql.DB
}
