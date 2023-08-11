package minetestconfig_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingParser(t *testing.T) {
	parts := strings.Split("my stuff  things", " ")
	assert.Equal(t, 3, len(parts))
}
