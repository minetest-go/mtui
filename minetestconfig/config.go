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
				st := sts[multiline_key]
				s.Add(multiline_key, multiline_value, st)

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

		value := strings.Trim(line[sepIndex+1:], " ")
		key := strings.Trim(line[:sepIndex], " ")

		if strings.TrimSpace(value) == multine_delimiter {
			in_multiline = true
			multiline_key = key
			continue
		}

		st := sts[key]
		s.Add(key, value, st)
	}
	return nil
}

func (s Settings) Write(w io.Writer, sts SettingTypes) error {
	if sts == nil {
		sts = SettingTypes{}
	}

	for key, setting := range s {
		entry := fmt.Sprintf("%s = ", key)
		value := setting.ToStringValue(sts[key])

		if strings.Contains(value, "\n") {
			// multiline value
			entry += multine_delimiter + "\n" + value + "\n" + multine_delimiter
		} else {
			// simple value
			entry += value + "\n"
		}
		_, err := w.Write([]byte(entry))
		if err != nil {
			return err
		}
	}
	return nil
}
