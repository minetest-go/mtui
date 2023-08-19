package depanalyzer_test

import (
	"mtui/minetestconfig/depanalyzer"
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

func TestParseModConfUnifiedInv(t *testing.T) {
	cfg, err := depanalyzer.ParseModConf([]byte(`
name = unified_inventory

optional_depends = default, creative, sfinv, datastorage
description = """
Unified Inventory replaces the default survival and creative inventory.
It adds a nicer interface and a number of features, such as a crafting guide.
"""
min_minetest_version = 5.4.0
	`))
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 0, len(cfg.Depends))
	assert.Equal(t, 4, len(cfg.OptionalDepends))
}
