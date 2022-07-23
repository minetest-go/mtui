package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"mtui/config"

	"github.com/stretchr/testify/assert"
)

func TestParseNonExisting(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "mtui.json")
	assert.NoError(t, err)

	_, err = config.Parse(f.Name())
	assert.Error(t, err)
}

func TestCreateParse(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "mtui.json")
	assert.NoError(t, err)

	cfg, err := config.WriteDefault(f.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "127.0.0.1", cfg.CookieDomain)
	assert.Equal(t, "/", cfg.CookiePath)
	assert.Equal(t, false, cfg.CookieSecure)
	assert.True(t, len(cfg.APIKey) > 0)
	assert.True(t, len(cfg.JWTKey) > 0)
}
