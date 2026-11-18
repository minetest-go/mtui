package bridge_test

import (
	"mtui/bridge"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	assert.Equal(t, "xy-z", bridge.SanitizeLuaString("xy-z"))
	assert.Equal(t, "xy-z", bridge.SanitizeLuaString("'xy-z"))
	assert.Equal(t, "xy-z", bridge.SanitizeLuaString("'xy-z''"))
	assert.Equal(t, "xy-z", bridge.SanitizeLuaString("\"'xy-z''"))
}

func TestQuoteLuaString(t *testing.T) {
	assert.Equal(t, `"xy-z"`, bridge.QuoteLuaString("xy-z"))
	assert.Equal(t, `"a\"b"`, bridge.QuoteLuaString(`a"b`))
	assert.Equal(t, `"a'b"`, bridge.QuoteLuaString(`a'b`))
	assert.Equal(t, `"a\\b"`, bridge.QuoteLuaString(`a\b`))
	assert.Equal(t, `"a\010b"`, bridge.QuoteLuaString("a\nb"))
	assert.Equal(t, `"\") os.exit() --"`, bridge.QuoteLuaString(`") os.exit() --`))
	assert.Equal(t, `"äö"`, bridge.QuoteLuaString("äö"))
}
