package modmanager_test

import (
	"mtui/api/cdb"
	"mtui/modmanager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	installed := []string{"default", "player_api"}
	cli := cdb.New()

	rd, err := modmanager.ResolveDependencies(cli, installed, "mt-mods", "pipeworks")
	assert.NoError(t, err)
	assert.NotNil(t, rd)
}
