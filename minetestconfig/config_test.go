package minetestconfig_test

import (
	"mtui/minetestconfig"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cfg := minetestconfig.MinetestConfig{}
	f, err := os.Open("testdata/simple.conf")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	err = cfg.Read(f)
	assert.NoError(t, err)

	assert.Equal(t, "123", cfg["xy"])
	assert.Equal(t, "a\nb\nc\n", cfg["description"])
	assert.Equal(t, 2, len(cfg))
}
