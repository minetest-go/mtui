package cdb_test

import (
	"mtui/api/cdb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResolveDependencies(t *testing.T) {
	c := cdb.New()
	cc := cdb.NewCachedClient(c, time.Hour*1)

	installed_pkgs := []string{"default"}
	selected_pkgs := []string{}

	rd, err := cdb.ResolveDependencies(cc, "mt-mods/technic_plus", selected_pkgs, installed_pkgs)
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	for _, di := range rd {
		switch di.Name {
		case "basic_materials":
			assert.True(t, len(di.Choices) >= 1)
		case "default":
			assert.True(t, di.Installed)
		}
	}

	assert.True(t, len(rd) >= 4)
}
