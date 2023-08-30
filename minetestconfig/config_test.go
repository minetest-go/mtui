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

myint = 456

mycoord = (1,2,3)
`

const simple_settingtypes = `
xy (my xa setting) string blah

myint (some int value) int 0

mycoord (some coord) v3f (0,0,0)
`

func TestParse(t *testing.T) {
	sts, err := minetestconfig.ParseSettingTypes([]byte(simple_settingtypes))
	assert.NoError(t, err)
	assert.NotNil(t, sts)

	cfg := minetestconfig.Settings{}
	buf := bytes.NewBuffer([]byte(simple_config))
	err = cfg.Read(buf, sts)
	assert.NoError(t, err)

	assert.Equal(t, 4, len(cfg))

	assert.NotNil(t, cfg["xy"])
	assert.Equal(t, "123", cfg["xy"].Value)

	assert.NotNil(t, cfg["description"])
	assert.Equal(t, "a\nb\nc\n", cfg["description"].Value)

	assert.NotNil(t, cfg["myint"])
	assert.Equal(t, "456", cfg["myint"].Value)

	assert.NotNil(t, sts["mycoord"])
	assert.NotNil(t, cfg["mycoord"])
	assert.Equal(t, 1.0, cfg["mycoord"].X)
}
