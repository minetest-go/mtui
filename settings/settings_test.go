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
	s := settings.New()
	err := s.Read(bytes.NewReader([]byte(txt)))
	assert.NoError(t, err)

	m := s.GetAll()
	assert.Equal(t, 3, len(m))
	assert.Equal(t, "1", m["x"])
	assert.Equal(t, "2", m["y"])
	assert.Equal(t, "(1,2,-3)", m["a"])

	buf := bytes.NewBuffer([]byte{})
	assert.NoError(t, s.Write(buf))

	s = settings.New()
	s.Read(buf)

	m = s.GetAll()
	assert.Equal(t, 3, len(m))
	assert.Equal(t, "1", m["x"])
	assert.Equal(t, "2", m["y"])
	assert.Equal(t, "(1,2,-3)", m["a"])
}
