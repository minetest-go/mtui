package bridge

import "strings"

// strips all string termination characters
func SanitizeLuaString(in string) string {
	in = strings.ReplaceAll(in, "'", "")
	in = strings.ReplaceAll(in, "\"", "")
	return in
}
