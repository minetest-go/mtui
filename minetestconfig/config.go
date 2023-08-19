package minetestconfig

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const multine_delimiter = "\"\"\""

func (s Settings) Read(r io.Reader) error {
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
				s[multiline_key] = multiline_value
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

		s[key] = value
	}
	return nil
}

func (s Settings) Write(w io.Writer) error {
	for key, value := range s {
		_, err := w.Write([]byte(fmt.Sprintf("%s = %s\n", key, value)))
		if err != nil {
			return err
		}
	}
	return nil
}
