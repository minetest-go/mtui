package minetestconfig

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"mtui/minetestconfig/depanalyzer"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var bracket_replacer = strings.NewReplacer(")", "", "(", "")

func ParseSettingTypes(data []byte) (SettingTypes, error) {
	sc := bufio.NewScanner(bytes.NewReader(data))

	stypes := SettingTypes{}

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
			Default:         &Setting{},
		}
		// reset comment for next entry
		last_comment = ""

		// disassemble setting line
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			// not a valid setting
			continue
		}
		s.Key = strings.TrimSpace(parts[0])

		// remove parsed key from line
		line = line[len(s.Key)+1:]

		i1 := strings.Index(line, "(")
		i2 := strings.Index(line, ")")

		if i1 < 0 || i2 < 0 {
			// invalid short desc
			continue
		}

		s.ShortDescription = line[i1+1 : i2]

		// remove parsed desc from line
		line = line[i2+2:]

		i1 = strings.Index(line, " ")
		if i1 < 0 {
			// end of line
			if line == "" {
				// empty type, skip
				continue
			} else {
				// last piece is type
				s.Type = line
			}
		} else {
			// line continues
			s.Type = line[:i1]
		}

		// remove parsed type from line
		line = line[i1+1:]
		s.Default.ParseStringValue(line, s)

		stypes[s.Key] = s
	}

	return stypes, nil
}

func GetAllSettingTypes(dir string) (SettingTypes, error) {
	sts := SettingTypes{}

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
				sts[s.Key] = s
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
			sts[s.Key] = s
		}

		return nil
	})

	return sts, err
}

//go:embed server_settings.txt
var serversettings embed.FS

func GetServerSettingTypes() (SettingTypes, error) {
	data, err := serversettings.ReadFile("server_settings.txt")
	if err != nil {
		return nil, err
	}

	return ParseSettingTypes(data)
}
