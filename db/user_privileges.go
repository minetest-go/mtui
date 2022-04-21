package db

import (
	"database/sql"
)

type UserPrivilegeRepository struct {
	db *sql.DB
}

type PrivilegeEntry struct {
	ID        int64  `json:"id"`
	Privilege string `json:"privilege"`
}

func NewUserPrivilegeRepository(db *sql.DB) *UserPrivilegeRepository {
	return &UserPrivilegeRepository{db: db}
}

func (repo *UserPrivilegeRepository) GetByID(id int64) ([]*PrivilegeEntry, error) {
	rows, err := repo.db.Query("select id,privilege from user_privileges where id = $1", id)
	if err != nil {
		return nil, err
	}
	list := make([]*PrivilegeEntry, 0)
	for rows.Next() {
		entry := &PrivilegeEntry{}
		err = rows.Scan(&entry.ID, &entry.Privilege)
		if err != nil {
			return nil, err
		}
		list = append(list, entry)
	}
	return list, nil
}

func (repo *UserPrivilegeRepository) Create(entry *PrivilegeEntry) error {
	_, err := repo.db.Exec("insert into user_privileges(id,privilege) values($1,$2)", entry.ID, entry.Privilege)
	return err
}

func (repo *UserPrivilegeRepository) Delete(id int64, privilege string) error {
	_, err := repo.db.Exec("delete from user_privileges where id = $1 and privilege = $2", id, privilege)
	return err
}
