package db

import (
	"mtui/types"

	"github.com/minetest-go/dbutil"
)

type Repositories struct {
	ModRepo              *ModRepository
	ConfigRepo           *ConfigRepository
	FeatureRepository    *FeatureRepository
	LogRepository        *LogRepository
	MetricTypeRepository *MetricTypeRepository
	MetricRepository     *MetricRepository
	OauthAppRepo         *OauthAppRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		ModRepo:              &ModRepository{dbu: dbutil.New[*types.Mod](db, dbutil.DialectSQLite, func() *types.Mod { return &types.Mod{} })},
		ConfigRepo:           &ConfigRepository{dbu: dbutil.New[*types.ConfigEntry](db, dbutil.DialectSQLite, func() *types.ConfigEntry { return &types.ConfigEntry{} })},
		FeatureRepository:    &FeatureRepository{dbu: dbutil.New[*types.Feature](db, dbutil.DialectSQLite, func() *types.Feature { return &types.Feature{} })},
		LogRepository:        &LogRepository{db: db, dbu: dbutil.New[*types.Log](db, dbutil.DialectSQLite, func() *types.Log { return &types.Log{} })},
		MetricTypeRepository: &MetricTypeRepository{dbu: dbutil.New[*types.MetricType](db, dbutil.DialectSQLite, func() *types.MetricType { return &types.MetricType{} })},
		MetricRepository:     &MetricRepository{dbu: dbutil.New[*types.Metric](db, dbutil.DialectSQLite, func() *types.Metric { return &types.Metric{} })},
		OauthAppRepo:         &OauthAppRepository{dbu: dbutil.New[*types.OauthApp](db, dbutil.DialectSQLite, func() *types.OauthApp { return &types.OauthApp{} })},
	}
}
