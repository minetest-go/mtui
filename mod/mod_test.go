package mod

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	installDir, err := os.CreateTemp(os.TempDir(), "mod")
	assert.NoError(t, err)
	assert.NotNil(t, installDir)

	// remove file
	err = os.Remove(installDir.Name())
	assert.NoError(t, err)

	// mkdir
	err = os.MkdirAll(installDir.Name(), 0755)
	assert.NoError(t, err)

	err = Install(installDir.Name())
	assert.NoError(t, err)
}
