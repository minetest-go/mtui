package modmanager_test

import (
	"mtui/modmanager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseModConf(t *testing.T) {

	cfg, err := modmanager.ParseModConf([]byte(`
		name = abc
		depends = xy, my_mod
	`))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "abc", cfg.Name)
	assert.Equal(t, 2, len(cfg.Depends))

	cfg, err = modmanager.ParseModConf([]byte(`
		name = abc
		depends = """
		xy,
		my_mod
		"""
	`))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "abc", cfg.Name)
	assert.Equal(t, 2, len(cfg.Depends))
}
