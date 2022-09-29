package db

import (
	"database/sql"
)

type Repositories struct {
	ModRepo    *ModRepository
	ConfigRepo *ConfigRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ModRepo:    &ModRepository{DB: db},
		ConfigRepo: &ConfigRepository{DB: db},
	}
}
