package db

import (
	"database/sql"
)

type Repositories struct {
	ModRepo              *ModRepository
	ConfigRepo           *ConfigRepository
	FeatureRepository    *FeatureRepository
	LogRepository        *LogRepository
	MetricTypeRepository *MetricTypeRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		ModRepo:              &ModRepository{DB: db},
		ConfigRepo:           &ConfigRepository{DB: db},
		FeatureRepository:    &FeatureRepository{DB: db},
		LogRepository:        &LogRepository{DB: db},
		MetricTypeRepository: &MetricTypeRepository{DB: db},
	}
}
