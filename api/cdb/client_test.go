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

	details, err := c.GetDetails(pkgs[0])
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, pkgs[0].Name, details.Name)

	pkgs, err = c.SearchPackages(&cdb.PackageQuery{Query: "blockexchange"})
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)
	assert.True(t, len(pkgs) > 0 && len(pkgs) < 100)

	deps, err := c.GetDependencies(pkgs[0])
	assert.NoError(t, err)
	assert.NotNil(t, deps)
	assert.True(t, len(deps) > 0)

	releases, err := c.GetReleases(pkgs[0])
	assert.NoError(t, err)
	assert.NotNil(t, releases)
	assert.True(t, len(releases) > 0)

	release, err := c.GetRelease(pkgs[0], releases[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, release)
	assert.NotEmpty(t, release.URL)

	screenshots, err := c.GetScreenshots(pkgs[0])
	assert.NoError(t, err)
	assert.NotNil(t, screenshots)
	assert.True(t, len(screenshots) > 0)
}
