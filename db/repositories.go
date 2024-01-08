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
	ChatLogRepo          *ChatLogRepository
	MetricTypeRepository *MetricTypeRepository
	MetricRepository     *MetricRepository
	OauthAppRepo         *OauthAppRepository
	MeseconsRepo         *MeseconsRepository
}

func NewRepositories(db dbutil.DBTx) *Repositories {
	return &Repositories{
		ModRepo:              &ModRepository{dbu: dbutil.New[*types.Mod](db, dbutil.DialectSQLite, func() *types.Mod { return &types.Mod{} })},
		ConfigRepo:           &ConfigRepository{dbu: dbutil.New[*types.ConfigEntry](db, dbutil.DialectSQLite, func() *types.ConfigEntry { return &types.ConfigEntry{} })},
		FeatureRepository:    &FeatureRepository{dbu: dbutil.New[*types.Feature](db, dbutil.DialectSQLite, func() *types.Feature { return &types.Feature{} })},
		LogRepository:        NewLogRepository(db),
		ChatLogRepo:          &ChatLogRepository{dbu: dbutil.New[*types.ChatLog](db, dbutil.DialectSQLite, func() *types.ChatLog { return &types.ChatLog{} })},
		MetricTypeRepository: &MetricTypeRepository{dbu: dbutil.New[*types.MetricType](db, dbutil.DialectSQLite, func() *types.MetricType { return &types.MetricType{} })},
		MetricRepository:     &MetricRepository{dbu: dbutil.New[*types.Metric](db, dbutil.DialectSQLite, func() *types.Metric { return &types.Metric{} })},
		OauthAppRepo:         &OauthAppRepository{dbu: dbutil.New[*types.OauthApp](db, dbutil.DialectSQLite, func() *types.OauthApp { return &types.OauthApp{} })},
		MeseconsRepo:         &MeseconsRepository{dbu: dbutil.New[*types.Mesecons](db, dbutil.DialectSQLite, func() *types.Mesecons { return &types.Mesecons{} })},
	}
}
