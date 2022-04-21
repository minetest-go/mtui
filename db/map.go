package db

import "database/sql"

const MAP_MIGRATE = `
pragma journal_mode = wal;
begin;
CREATE TABLE if not exists blocks (
	pos INT PRIMARY KEY,
	data BLOB
);
commit;
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

func (repo *MapRepository) GetSize() (int64, error) {
	res, err := repo.db.Query("SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size()")
	if err != nil {
		return 0, err
	}
	var size int64
	res.Next()
	err = res.Scan(&size)
	return size, err
}

func (repo *MapRepository) CountBlocks() (int64, error) {
	res, err := repo.db.Query("SELECT count(*) from blocks")
	if err != nil {
		return 0, err
	}
	var count int64
	res.Next()
	err = res.Scan(&count)
	return count, err
}
