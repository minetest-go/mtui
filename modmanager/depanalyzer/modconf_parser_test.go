package depanalyzer_test

import (
	"mtui/modmanager/depanalyzer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseModConf(t *testing.T) {

	cfg, err := depanalyzer.ParseModConf([]byte(`
		name = abc
		depends = xy, my_mod
	`))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "abc", cfg.Name)
	assert.Equal(t, 2, len(cfg.Depends))

	cfg, err = depanalyzer.ParseModConf([]byte(`
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

	cfg, err = depanalyzer.ParseModConf([]byte(`
		name = abc
		description = """
		something, something
		something else
		"""
	`))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEqual(t, "", cfg.Description)
}
