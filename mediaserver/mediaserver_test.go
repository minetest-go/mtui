package mediaserver_test

import (
	"mtui/mediaserver"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaServer(t *testing.T) {
	m := mediaserver.New()
	wd, err := os.Getwd()
	dir := path.Join(wd, "testdata")
	assert.NoError(t, err)
	m.Scan(dir)

	assert.Equal(t, 2, m.Count)
	assert.Equal(t, int64(9), m.Size)
	assert.NotEqual(t, "", m.Media["7110eda4d09e062aa5e4a390b0a572ac0d2c0220"])
	assert.NotEqual(t, "", m.Media["8cb2237d0679ca88db6464eac60da96345513964"])
}
