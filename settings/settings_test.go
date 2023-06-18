package settings_test

import (
	"bytes"
	"mtui/settings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const txt = `
x = 1
y=2
#comment
a = (1,2,-3)
`

func TestSettings(t *testing.T) {
	s := settings.Settings{}
	err := s.Read(bytes.NewReader([]byte(txt)))
	assert.NoError(t, err)

	assert.Equal(t, 3, len(s))
	assert.Equal(t, "1", s["x"])
	assert.Equal(t, "2", s["y"])
	assert.Equal(t, "(1,2,-3)", s["a"])

	buf := bytes.NewBuffer([]byte{})
	assert.NoError(t, s.Write(buf))

	s = settings.Settings{}
	s.Read(buf)

	assert.Equal(t, 3, len(s))
	assert.Equal(t, "1", s["x"])
	assert.Equal(t, "2", s["y"])
	assert.Equal(t, "(1,2,-3)", s["a"])
}
