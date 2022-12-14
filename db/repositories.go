package db

import (
	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	ModRepo              *ModRepository
	ConfigRepo           *ConfigRepository
	FeatureRepository    *FeatureRepository
	LogRepository        *LogRepository
	MetricTypeRepository *MetricTypeRepository
	MetricRepository     *MetricRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		ModRepo:              &ModRepository{DB: db},
		ConfigRepo:           &ConfigRepository{DB: db},
		FeatureRepository:    &FeatureRepository{DB: db},
		LogRepository:        &LogRepository{DB: db},
		MetricTypeRepository: &MetricTypeRepository{DB: db},
		MetricRepository:     &MetricRepository{DB: db},
	}
}
