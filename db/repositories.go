package db

import (
	"database/sql"
)

type Repositories struct {
	ModRepo *ModRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ModRepo: &ModRepository{DB: db},
	}
}
