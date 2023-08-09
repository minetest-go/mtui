package depanalyzer

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type ModConfig struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Depends         []string `json:"depends"`
	OptionalDepends []string `json:"optional_depends"`
	Author          string   `json:"author"`
	Title           string   `json:"title"`
}

func addModConfField(cfg *ModConfig, key, value string) {
	switch key {
	case "name":
		cfg.Name = value
	case "description":
		cfg.Description = value
	case "author":
		cfg.Author = value
	case "title":
		cfg.Title = value
	case "depends":
		value = strings.ReplaceAll(value, "\n", "")
		value = strings.ReplaceAll(value, "\r", "")

		for _, dep := range strings.Split(value, ",") {
			cfg.Depends = append(cfg.Depends, strings.TrimSpace(dep))
		}
	case "optional_depends":
		value = strings.ReplaceAll(value, "\n", "")
		value = strings.ReplaceAll(value, "\r", "")

		for _, dep := range strings.Split(value, ",") {
			cfg.OptionalDepends = append(cfg.OptionalDepends, strings.TrimSpace(dep))
		}
	}
}

func ParseModConf(data []byte) (*ModConfig, error) {
	cfg := &ModConfig{
		Depends:         make([]string, 0),
		OptionalDepends: make([]string, 0),
	}
	sc := bufio.NewScanner(bytes.NewReader(data))

	is_in_multiline := false
	multiline_key := ""
	multiline_value := ""

	line_num := 0
	for sc.Scan() {
		line_num++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		if !is_in_multiline {
			if strings.HasPrefix(line, "#") {
				// comment
				continue
			}
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}

			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])

			if strings.Contains(value, "\"\"\"") {
				is_in_multiline = true
				multiline_key = key
			} else {
				addModConfField(cfg, key, value)
			}
		} else {
			// multiline case
			if strings.TrimSpace(line) == "\"\"\"" {
				//end of multiline
				is_in_multiline = false

				addModConfField(cfg, multiline_key, multiline_value)
				multiline_value = ""
			} else {
				// append all found lines
				multiline_value += line
			}
		}

	}

	if is_in_multiline {
		return nil, fmt.Errorf("unterminated multiline string encountered")
	}

	return cfg, nil
}
