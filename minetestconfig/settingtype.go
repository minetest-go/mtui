package minetestconfig

// setting with type information
type SettingType struct {
	Key              string   `json:"key"`   // my.key
	Category         []string `json:"topic"` // stacked list of categories ["Advanced", "Developer Options", "Mod Security"]
	LongDescription  string   `json:"long_description"`
	ShortDescription string   `json:"short_description"`
	Type             string   `json:"type"`    // int, string, float, bool, enum
	Choices          []string `json:"choices"` // enum choices
	Default          *Setting `json:"default"`
	Min              float64  `json:"min"`
	Max              float64  `json:"max"`
	DefaultMGFlags   []string `json:"default_mg_flags"`
}
