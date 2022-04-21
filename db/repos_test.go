package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepositories(t *testing.T) {
	repos, err := CreateRepositories(os.TempDir())
	assert.NoError(t, err)
	assert.NotNil(t, repos)
}
