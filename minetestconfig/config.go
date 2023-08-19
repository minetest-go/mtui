package minetestconfig

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const multine_delimiter = "\"\"\""

func (s Settings) Read(r io.Reader, sts SettingTypes) error {
	if sts == nil {
		//default to empty settingtypes
		sts = SettingTypes{}
	}
	sc := bufio.NewScanner(r)
	linenum := 0
	in_multiline := false
	multiline_value := ""
	multiline_key := ""

	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		linenum++

		if in_multiline {
			if line == multine_delimiter {
				// end of multiline
				in_multiline = false
				setting := s[multiline_key]
				if setting == nil {
					setting = &Setting{}
					s[multiline_key] = setting
				}

				setting.Value = multiline_value
			} else {
				// continue with multiline
				multiline_value += line + "\n"
			}
			continue
		}

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		sepIndex := strings.Index(line, "=")
		if sepIndex < 0 {
			continue
		}
		//TODO: multiline values with """

		value := strings.Trim(line[sepIndex+1:], " ")
		key := strings.Trim(line[:sepIndex], " ")

		if strings.TrimSpace(value) == multine_delimiter {
			in_multiline = true
			multiline_key = key
			continue
		}

		setting := s[multiline_key]
		if setting == nil {
			setting = &Setting{}
			s[key] = setting
		}

		setting.Value = value
	}
	return nil
}

func (s Settings) Write(w io.Writer) error {
	for key, setting := range s {
		entry := fmt.Sprintf("%s = ", key)

		if strings.Contains(setting.Value, "\n") {
			// multiline value
			entry += multine_delimiter + "\n" + setting.Value + "\n" + multine_delimiter
		} else {
			// simple value
			entry += setting.Value + "\n"
		}
		_, err := w.Write([]byte(entry))
		if err != nil {
			return err
		}
	}
	return nil
}
