package minetestconfig

import (
	"bufio"
	"bytes"
	"fmt"
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
			LongDescription: last_comment,
			Category:        categories,
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
		if len(parts) < 2 {
			// not enough parts
			continue
		}
		s.Type = strings.TrimSpace(parts[0])
		s.Default = strings.TrimSpace(parts[1])

		if len(parts) >= 3 {
			// float 600.0 0.0
			v, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid 'min' setting in '%s': %v", s.Key, err)
			}
			s.Min = v
		}

		if len(parts) >= 4 {
			// int 20 -1 32767
			v, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid 'max' setting in '%s': %v", s.Key, err)
			}
			s.Max = v
		}

		list = append(list, s)
	}

	return list, nil
}
