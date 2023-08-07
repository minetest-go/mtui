package cdb_test

import (
	"mtui/api/cdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackages(t *testing.T) {
	c := cdb.New()
	pkgs, err := c.GetPackages()
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)
	assert.True(t, len(pkgs) > 0)
}
