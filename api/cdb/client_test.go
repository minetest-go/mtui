package cdb_test

import (
	"mtui/api/cdb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackages(t *testing.T) {
	c := cdb.New()

	tags, err := c.GetTags()
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.True(t, len(tags) > 0)

	cws, err := c.GetContentWarnings()
	assert.NoError(t, err)
	assert.NotNil(t, cws)
	assert.True(t, len(cws) > 0)

	pkgs, err := c.SearchPackages(&cdb.PackageQuery{Query: "blockexchange", Author: "buckaroobanzay", Limit: 1})
	assert.NoError(t, err)
	assert.NotNil(t, pkgs)
	assert.Equal(t, 1, len(pkgs))

	details, err := c.GetDetails(pkgs[0].Author, pkgs[0].Name)
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, pkgs[0].Name, details.Name)

	deps, err := c.GetDependencies(pkgs[0].Author, pkgs[0].Name)
	assert.NoError(t, err)
	assert.NotNil(t, deps)
	assert.True(t, len(deps) > 0)

	releases, err := c.GetReleases(pkgs[0].Author, pkgs[0].Name)
	assert.NoError(t, err)
	assert.NotNil(t, releases)
	assert.True(t, len(releases) > 0)

	release, err := c.GetRelease(pkgs[0].Author, pkgs[0].Name, releases[0].ID)
	assert.NoError(t, err)
	assert.NotNil(t, release)
	assert.NotEmpty(t, release.URL)

	z, err := c.DownloadZip(releases[0])
	assert.NoError(t, err)
	assert.NotNil(t, z)
	f, err := z.Open("blockexchange/init.lua")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	screenshots, err := c.GetScreenshots(pkgs[0].Author, pkgs[0].Name)
	assert.NoError(t, err)
	assert.NotNil(t, screenshots)
	assert.True(t, len(screenshots) > 0)
}
