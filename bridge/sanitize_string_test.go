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
