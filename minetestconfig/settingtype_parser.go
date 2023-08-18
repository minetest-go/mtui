package minetestconfig

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"mtui/modmanager/depanalyzer"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func ParseSettingTypes(data []byte) ([]*SettingType, error) {
	sc := bufio.NewScanner(bytes.NewReader(data))

	list := make([]*SettingType, 0)

	last_comment := ""
	categories := []string{}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line == "" {
			// empty line
			continue
		}

		if strings.HasPrefix(line, "#") {
			// append comment
			last_comment = fmt.Sprintf("%s\n%s", last_comment, strings.TrimSpace(strings.TrimPrefix(line, "#")))
			continue
		}

		if strings.HasPrefix(line, "[") {
			// category
			category_depth := strings.Count(line, "*")
			category := strings.NewReplacer("[", "", "]", "", "*", "").Replace(line)
			if len(categories) > category_depth {
				// strip outer categories
				categories = categories[0:category_depth]
			}
			categories = append(categories, category)
			continue
		}

		// everything else is a settingtype entry
		s := &SettingType{
			LongDescription: strings.TrimPrefix(last_comment, "\n"),
			Category:        append([]string{}, categories...),
		}
		// reset comment for next entry
		last_comment = ""

		// disassemble setting line
		parts := strings.Split(line, "(")
		if len(parts) < 2 {
			// "(" not found, skip
			continue
		}
		s.Key = strings.TrimSpace(parts[0])

		descparts := strings.Split(parts[1], ")")
		if len(descparts) < 2 {
			// ")" not found
			continue
		}
		s.ShortDescription = strings.TrimSpace(descparts[0])

		rest := strings.TrimSpace(descparts[1])
		// remove double spaces
		for strings.Contains(rest, "  ") {
			rest = strings.ReplaceAll(rest, "  ", " ")
		}

		parts = strings.Split(rest, " ")

		if len(parts) >= 1 {
			s.Type = strings.TrimSpace(parts[0])
		}

		if len(parts) >= 2 {

			switch s.Type {
			case "v3f":
				// mgfractal_scale (Scale) v3f (4096.0, 1024.0, 4096.0)
				v3f_str := strings.Join(parts[1:], "")
				v3f_str = strings.ReplaceAll(v3f_str, "(", "")
				v3f_str = strings.ReplaceAll(v3f_str, "(", "")

				v3f := strings.Split(v3f_str, ",")
				if len(v3f) != 3 {
					break
				}

				s.X, _ = strconv.ParseFloat(v3f[0], 64)
				s.Y, _ = strconv.ParseFloat(v3f[1], 64)
				s.Z, _ = strconv.ParseFloat(v3f[2], 64)
			default:
				s.Default = strings.TrimSpace(parts[1])
			}
		}

		if len(parts) >= 3 {
			switch s.Type {
			case "string":
				// server_name (Server name) string Minetest server
				s.Default = strings.Join(parts[2:], " ")
			case "enum", "flags":
				// hudbars_bar_type (HUD bars style) enum progress_bar progress_bar,statbar_classic,statbar_modern
				// mg_flags (Mapgen flags) flags caves,dungeons,light,decorations,biomes,ores caves,dungeons,light,decorations,biomes,ores,nocaves,nodungeons,nolight,nodecorations,nobiomes,noores
				s.Choices = strings.Split(parts[2], ",")
			case "int", "float":
				// float 600.0 0.0
				v, err := strconv.ParseFloat(parts[2], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid 'min' setting in '%s': %v", s.Key, err)
				}
				s.Min = v
			}
		}

		if len(parts) >= 4 {
			// int 20 -1 32767
			switch s.Type {
			case "int", "float":
				v, err := strconv.ParseFloat(parts[3], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid 'max' setting in '%s': %v", s.Key, err)
				}
				s.Max = v
			}
		}

		list = append(list, s)
	}

	return list, nil
}

func GetAllSettingTypes(dir string) ([]*SettingType, error) {
	list := []*SettingType{}

	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, _ error) error {
		if d != nil && d.IsDir() {
			return nil
		}

		basename := path.Base(p)
		if basename != "settingtypes.txt" {
			return nil
		}

		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		st, err := ParseSettingTypes(data)
		if err != nil {
			return err
		}

		dirname := path.Dir(p)
		gameconf, _ := os.Stat(path.Join(dirname, "game.conf"))
		if gameconf != nil {
			// game-setting
			for _, s := range st {
				s.Category = append([]string{"Game"}, s.Category...)
				list = append(list, s)
			}
			return nil
		}

		modname := path.Dir(dirname)
		modconf, _ := os.Stat(path.Join(dirname, "mod.conf"))
		if modconf != nil {
			// mod-setting
			data, err = os.ReadFile(path.Join(dirname, "mod.conf"))
			if err != nil {
				return err
			}
			c, err := depanalyzer.ParseModConf(data)
			if err != nil {
				return err
			}
			modname = c.Name
		}
		dependstxt, _ := os.Stat(path.Join(dirname, "depends.txt"))
		if dependstxt == nil && modconf == nil {
			// not a mod directory
			return nil
		}

		for _, s := range st {
			s.Category = append([]string{"Mods", modname}, s.Category...)
			list = append(list, s)
		}

		return nil
	})

	return list, err
}

//go:embed server_settings.txt
var serversettings embed.FS

func GetServerSettingTypes() ([]*SettingType, error) {
	data, err := serversettings.ReadFile("server_settings.txt")
	if err != nil {
		return nil, err
	}

	return ParseSettingTypes(data)
}
