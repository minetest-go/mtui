package minetestconfig

type MinetestConfig map[string]string

type SettingType struct {
	LongDescription  string `json:"long_description"`
	ShortDescription string `json:"short_description"`
	Key              string `json:"key"`
	Type             string `json:"type"`
	Default          string `json:"default"`
}
