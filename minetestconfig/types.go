package minetestconfig

type Settings map[string]*Setting

func (ss Settings) Add(key, value string, st *SettingType) {
	s := &Setting{}

	if st == nil {
		// use string type as default
		st = &SettingType{Type: "string"}
	}

	s.ParseStringValue(value, st)
	ss[key] = s
}

type SettingTypes map[string]*SettingType
