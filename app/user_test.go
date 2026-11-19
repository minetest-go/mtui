package app_test

import (
	"mtui/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	a := testutil.NewApp(t)
	e, err := a.CreateAdmin("admin", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, e)
}
