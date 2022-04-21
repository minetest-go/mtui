package minetestconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cfg, err := Parse("testdata/simple.conf")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "123", cfg["xy"])
}
