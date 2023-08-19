package minetestconfig_test

import (
	"bytes"
	"mtui/minetestconfig"
	"testing"

	"github.com/stretchr/testify/assert"
)

const simple_config = `
xy = 123
# abc = def

description = """
a
b
c
"""
`

const simple_settingtypes = `
xy (my xa setting) string blah
`

func TestParse(t *testing.T) {
	st, err := minetestconfig.ParseSettingTypes([]byte(simple_settingtypes))
	assert.NoError(t, err)
	assert.NotNil(t, st)
	//TODO: use settingtypes for config parsing

	cfg := minetestconfig.Settings{}
	buf := bytes.NewBuffer([]byte(simple_config))
	err = cfg.Read(buf)
	assert.NoError(t, err)

	assert.Equal(t, "123", cfg["xy"])
	assert.Equal(t, "a\nb\nc\n", cfg["description"])
	assert.Equal(t, 2, len(cfg))
}
