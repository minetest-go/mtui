package db

import (
	"database/sql"
)

type Repositories struct {
	Mod *ModRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Mod: &ModRepository{db: db},
	}
}
