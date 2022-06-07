package modmanager_test

import (
	"mtui/modmanager"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	app := CreateTestApp(t)
	mm := modmanager.New(app.WorldDir)
	assert.NotNil(t, mm)

	mod := &types.Mod{
		Name: "",
	}

	s, err := mm.IsSync(mod)
	assert.NoError(t, err)
	assert.False(t, s)

	err = mm.Sync(mod)
	assert.NoError(t, err)

	// TODO
}
