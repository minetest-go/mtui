package settings

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Settings map[string]string

func (s Settings) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	linenum := 0
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		linenum++

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.Split(line, "=")

		if len(parts) != 2 {
			return fmt.Errorf("invalid delimiter count in line %d", linenum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

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
