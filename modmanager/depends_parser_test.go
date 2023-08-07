package modmanager_test

import (
	"mtui/modmanager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDependsTXT(t *testing.T) {
	deps := modmanager.ParseDependsTXT([]byte("xy, abc,           e"))
	assert.Equal(t, 3, len(deps))
	assert.Equal(t, "xy", deps[0])
	assert.Equal(t, "abc", deps[1])
	assert.Equal(t, "e", deps[2])
}
