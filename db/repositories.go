package db

import (
	"database/sql"
)

type Repositories struct {
	ModRepo           *ModRepository
	ConfigRepo        *ConfigRepository
	FeatureRepository *FeatureRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ModRepo:           &ModRepository{DB: db},
		ConfigRepo:        &ConfigRepository{DB: db},
		FeatureRepository: &FeatureRepository{DB: db},
	}
}
