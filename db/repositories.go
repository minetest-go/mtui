package db

import (
	"github.com/minetest-go/dbutil"
	"gorm.io/gorm"
)

type Repositories struct {
	ModRepo              *ModRepository
	ConfigRepo           *ConfigRepository
	FeatureRepository    *FeatureRepository
	LogRepository        *LogRepository
	ChatLogRepo          *ChatLogRepository
	MetricTypeRepository *MetricTypeRepository
	MetricRepository     *MetricRepository
	OauthAppRepo         *OauthAppRepository
	MeseconsRepo         *MeseconsRepository
}

func NewRepositories(db dbutil.DBTx, g *gorm.DB) *Repositories {
	return &Repositories{
		ModRepo:              &ModRepository{g: g},
		ConfigRepo:           &ConfigRepository{g: g},
		FeatureRepository:    &FeatureRepository{g: g},
		LogRepository:        &LogRepository{g: g},
		ChatLogRepo:          &ChatLogRepository{g: g},
		MetricTypeRepository: &MetricTypeRepository{g: g},
		MetricRepository:     &MetricRepository{g: g},
		OauthAppRepo:         &OauthAppRepository{g: g},
		MeseconsRepo:         &MeseconsRepository{g: g},
	}
}
