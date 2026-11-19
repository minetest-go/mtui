package web_test

import (
	"mtui/testutil"
	"mtui/web"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebSetup(t *testing.T) {
	app := testutil.NewApp(t)
	assert.NoError(t, web.Setup(app))
}
