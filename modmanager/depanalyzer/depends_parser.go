package depanalyzer

import "strings"

func ParseDependsTXT(data []byte) []string {
	deps := []string{}
	for _, dep := range strings.Split(string(data), ",") {
		deps = append(deps, strings.TrimSpace(dep))
	}
	return deps
}
