package matterbridge_test

import (
	"mtui/service/matterbridge"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data, err := os.ReadFile("testdata/matterbridge.simple.toml")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	cfg, err := matterbridge.ParseConfig(data)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// TODO
}
