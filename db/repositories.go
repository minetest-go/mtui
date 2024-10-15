package db

import (
	"mtui/types"

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
		MetricTypeRepository: &MetricTypeRepository{dbu: dbutil.New[*types.MetricType](db, dbutil.DialectSQLite, func() *types.MetricType { return &types.MetricType{} })},
		MetricRepository:     &MetricRepository{dbu: dbutil.New[*types.Metric](db, dbutil.DialectSQLite, func() *types.Metric { return &types.Metric{} })},
		OauthAppRepo:         &OauthAppRepository{dbu: dbutil.New[*types.OauthApp](db, dbutil.DialectSQLite, func() *types.OauthApp { return &types.OauthApp{} })},
		MeseconsRepo:         &MeseconsRepository{dbu: dbutil.New[*types.Mesecons](db, dbutil.DialectSQLite, func() *types.Mesecons { return &types.Mesecons{} })},
	}
}
