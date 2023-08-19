package minetestconfig

type Settings map[string]string

type SettingTypes map[string]*SettingType

type Setting struct {
	Key         string  `json:"key"`    // my.key
	Value       string  `json:"value"`  // value for primitive types
	X           float64 `json:"x"`      // v3f
	Y           float64 `json:"y"`      // v3f
	Z           float64 `json:"z"`      // v3f
	Offset      float64 `json:"offset"` // noise_params_*d
	Scale       float64 `json:"scale"`
	SpreadX     float64 `json:"spread_x"`
	SpreadY     float64 `json:"spread_y"`
	SpreadZ     float64 `json:"spread_z"`
	Seed        string  `json:"seed"`
	Octaves     float64 `json:"octaves"`
	Persistence float64 `json:"persistence"`
	Lacunarity  float64 `json:"lacunarity"`
}

type SettingType struct {
	*Setting
	Category         []string `json:"topic"` // stacked list of categories ["Advanced", "Developer Options", "Mod Security"]
	LongDescription  string   `json:"long_description"`
	ShortDescription string   `json:"short_description"`
	Type             string   `json:"type"`    // int, string, float, bool, enum
	Choices          []string `json:"choices"` // enum choices
	Default          string   `json:"default"`
	Min              float64  `json:"min"`
	Max              float64  `json:"max"`
	DefaultMGFlags   []string `json:"default_mg_flags"`
}
