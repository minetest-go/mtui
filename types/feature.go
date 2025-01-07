package types

const (
	FEATURE_API             = "api"
	FEATURE_MEDIASERVER     = "mediaserver"
	FEATURE_MAIL            = "mail"
	FEATURE_AREAS           = "areas"
	FEATURE_SKINSDB         = "skinsdb"
	FEATURE_LUASHELL        = "luashell"
	FEATURE_SHELL           = "shell"
	FEATURE_MODMANAGEMENT   = "modmanagement"
	FEATURE_XBAN            = "xban"
	FEATURE_MINETEST_CONFIG = "minetest_config"
	FEATURE_OTP             = "otp"
	FEATURE_DOCKER          = "docker"
	FEATURE_SIGNUP          = "signup"
	FEATURE_MESECONS        = "mesecons"
	FEATURE_ATM             = "atm"
	FEATURE_MINETEST_WEB    = "minetest_web"
)

type Feature struct {
	Name    string `json:"name" gorm:"primarykey;column:name"`
	Enabled bool   `json:"enabled" gorm:"column:enabled"`
}

func (m *Feature) TableName() string {
	return "feature"
}
