package types

type FeatureName string

const (
	FEATURE_API             FeatureName = "api"
	FEATURE_MEDIASERVER     FeatureName = "mediaserver"
	FEATURE_MAIL            FeatureName = "mail"
	FEATURE_AREAS           FeatureName = "areas"
	FEATURE_SKINSDB         FeatureName = "skinsdb"
	FEATURE_LUASHELL        FeatureName = "luashell"
	FEATURE_SHELL           FeatureName = "shell"
	FEATURE_MODMANAGEMENT   FeatureName = "modmanagement"
	FEATURE_XBAN            FeatureName = "xban"
	FEATURE_MINETEST_CONFIG FeatureName = "minetest_config"
	FEATURE_OTP             FeatureName = "otp"
	FEATURE_DOCKER          FeatureName = "docker"
	FEATURE_SIGNUP          FeatureName = "signup"
	FEATURE_MESECONS        FeatureName = "mesecons"
	FEATURE_ATM             FeatureName = "atm"
	FEATURE_MINETEST_WEB    FeatureName = "minetest_web"
)

type Feature struct {
	Name    FeatureName `json:"name" gorm:"primarykey;column:name"`
	Enabled bool        `json:"enabled" gorm:"column:enabled"`
}

func (m *Feature) TableName() string {
	return "feature"
}
