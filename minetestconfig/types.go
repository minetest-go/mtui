package minetestconfig

type MinetestConfig map[string]string

type SettingType struct {
	Category         []string `json:"topic"` // stacked list of categories ["Advanced", "Developer Options", "Mod Security"]
	LongDescription  string   `json:"long_description"`
	ShortDescription string   `json:"short_description"`
	Key              string   `json:"key"`     // my.key
	Type             string   `json:"type"`    // int, string, float, bool, enum
	Choices          []string `json:"choices"` // enum choices
	Default          string   `json:"default"`
	Min              float64  `json:"min"`
	Max              float64  `json:"max"`
	X                float64  `json:"x"`
	Y                float64  `json:"y"`
	Z                float64  `json:"z"`
}
