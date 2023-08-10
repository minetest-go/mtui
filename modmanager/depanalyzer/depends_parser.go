package depanalyzer

import "strings"

func ParseDependsTXT(data []byte) (*DependsInfo, error) {
	di := &DependsInfo{
		Depends:         make([]string, 0),
		OptionalDepends: make([]string, 0),
	}

	str := string(data)
	str = strings.ReplaceAll(str, "\n", ",")

	for _, raw_dep := range strings.Split(str, ",") {
		dep := strings.TrimSpace(raw_dep)
		if dep == "" {
			continue
		}
		if strings.HasSuffix(dep, "?") {
			// optional
			di.OptionalDepends = append(di.OptionalDepends, strings.ReplaceAll(dep, "?", ""))
		} else {
			// hard
			di.Depends = append(di.Depends, dep)
		}
	}

	return di, nil
}
