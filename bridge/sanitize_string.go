package bridge

import (
	"fmt"
	"strings"
)

// strips all string termination characters
func SanitizeLuaString(in string) string {
	in = strings.ReplaceAll(in, "'", "")
	in = strings.ReplaceAll(in, "\"", "")
	return in
}

// QuoteLuaString returns a double-quoted lua string literal for the input,
// escaping quotes, backslashes and control characters (lua 5.1 compatible)
func QuoteLuaString(in string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for i := 0; i < len(in); i++ {
		b := in[i]
		switch {
		case b == '"' || b == '\\':
			sb.WriteByte('\\')
			sb.WriteByte(b)
		case b < 0x20 || b == 0x7f:
			fmt.Fprintf(&sb, "\\%03d", b)
		default:
			sb.WriteByte(b)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}
