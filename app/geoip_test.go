package app_test

import (
	"io/ioutil"
	"mtui/app"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func copy(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, 0644)
}

func TestGeoipResolve(t *testing.T) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "mtui")
	assert.NoError(t, err)

	err = copy("testdata/GeoLite2-ASN-Test.mmdb", path.Join(tmpdir, "GeoLite2-ASN.mmdb"))
	assert.NoError(t, err)

	err = copy("testdata/GeoLite2-City-Test.mmdb", path.Join(tmpdir, "GeoLite2-City.mmdb"))
	assert.NoError(t, err)

	resolver := app.NewGeoipResolver(tmpdir)
	assert.NotNil(t, resolver)

	result := resolver.Resolve("81.2.69.142")
	assert.NotNil(t, result)

	assert.Equal(t, uint(0), result.ASN)
	assert.Equal(t, "London", result.City)
	assert.Equal(t, "United Kingdom", result.Country)
}
