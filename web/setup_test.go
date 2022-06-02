package web_test

import (
	"mtui/web"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebSetup(t *testing.T) {
	app := CreateTestApp(t)
	assert.NoError(t, web.Setup(app))
}
