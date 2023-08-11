package minetestconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func ParseSettingTypes(data []byte) ([]*SettingType, error) {
	sc := bufio.NewScanner(bytes.NewReader(data))

	list := make([]*SettingType, 0)

	var previous_comment string

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line == "" {
			// empty line
			continue
		}

		if strings.HasPrefix(line, "#") {
			// remember prefix and continue
			previous_comment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}

		// remove repeated spaces
		for strings.Contains(line, "  ") {
			line = strings.ReplaceAll(line, "  ", " ")
		}

		parts := strings.Split(line, " ")

		fmt.Println(previous_comment, parts)
	}

	return list, nil
}
