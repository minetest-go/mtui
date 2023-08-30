package minetestconfig

import (
	"fmt"
	"strconv"
	"strings"
)

// "bare" setting without type information
type Setting struct {
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

func (s *Setting) ParseStringValue(value string, st *SettingType) {
	switch st.Type {
	case "string", "path", "filepath", "key":
		// server_name (Server name) string Minetest server
		s.Value = value
	case "bool":
		parts := strings.Split(value, " ")
		if len(parts) >= 1 {
			s.Value = parts[0]
		}
	case "int", "float":
		parts := strings.Split(value, " ")
		if len(parts) >= 1 {
			s.Value = parts[0]
		}
		if len(parts) >= 2 {
			min, _ := strconv.ParseFloat(parts[1], 64)
			st.Min = &min
		}
		if len(parts) >= 3 {
			max, _ := strconv.ParseFloat(parts[2], 64)
			st.Max = &max
		}
	case "enum", "flags":
		parts := strings.Split(value, " ")
		if len(parts) >= 1 {
			s.Value = parts[0]
		}
		if len(parts) >= 2 {
			st.Choices = strings.Split(parts[1], ",")
		}
	case "v3f":
		i1 := strings.Index(value, "(")
		i2 := strings.Index(value, ")")
		if i1 < 0 || i2 < 0 {
			return
		}
		parts := strings.Split(value[i1+1:i2-1], ",")
		if len(parts) != 3 {
			return
		}
		s.X, _ = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		s.Y, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		s.Z, _ = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	case "noise_params_2d", "noise_params_3d":
		// <offset>, <scale>, (<spreadX>, <spreadY>, <spreadZ>), <seed>, <octaves>, <persistence>, <lacunarity>[, <default flags>]
		// mgfractal_np_seabed (Seabed noise) noise_params_2d -14, 9, (600, 600, 600), 41900, 5, 0.6, 2.0, eased
		// mgv5_np_cave1 (Cave1 noise) noise_params_3d 0, 12, (61, 61, 61), 52534, 3, 0.5, 2.0

		// remove brackets
		value = bracket_replacer.Replace(value)
		parts := strings.Split(value, ",")
		if len(parts) >= 9 {
			s.Offset, _ = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			s.Scale, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			s.SpreadX, _ = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			s.SpreadY, _ = strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
			s.SpreadZ, _ = strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
			s.Seed = strings.TrimSpace(parts[5])
			s.Octaves, _ = strconv.ParseFloat(strings.TrimSpace(parts[6]), 64)
			s.Persistence, _ = strconv.ParseFloat(strings.TrimSpace(parts[7]), 64)
			s.Lacunarity, _ = strconv.ParseFloat(strings.TrimSpace(parts[8]), 64)
		}

		if len(parts) >= 10 {
			st.DefaultMGFlags = strings.Split(strings.TrimSpace(parts[9]), ",")
		}
	}
}

func (s *Setting) ToStringValue(st *SettingType) string {
	if st == nil {
		return s.Value
	}

	switch st.Type {
	case "string", "path", "filepath", "key", "bool", "int", "float", "enum", "flags":
		return s.Value
	case "v3f":
		return fmt.Sprintf("(%.2f,%.2f,%.2f)", s.X, s.Y, s.Z)
	case "noise_params_2d", "noise_params_3d":
		return fmt.Sprintf("%.2f,%.2f,(%.2f,%.2f,%.2f), %s, %.2f, %.2f, %.2f",
			s.Offset, s.Scale, s.SpreadX, s.SpreadY, s.SpreadZ, s.Seed, s.Octaves, s.Persistence, s.Lacunarity)
	}
	return ""
}
