package depanalyzer_test

import (
	"mtui/minetestconfig/depanalyzer"
	"testing"

	"github.com/stretchr/testify/assert"
)

var villages_depends = `
handle_schematics
default
doors?
farming?
wool?
stairs?
cottages?
moretrees?
trees?
forest?
dryplants?
cavestuff?
snow?
moresnow?
darkage?
ethereal?
deco?
metals?
grounds?
moreblocks?
bell?
mg?
intllib?
mob_world_interaction?
cavestuff?
`

func TestParseDependsTXT(t *testing.T) {
	di, err := depanalyzer.ParseDependsTXT([]byte("xy, abc,           e"))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(di.Depends))
	assert.Equal(t, "xy", di.Depends[0])
	assert.Equal(t, "abc", di.Depends[1])
	assert.Equal(t, "e", di.Depends[2])

	di, err = depanalyzer.ParseDependsTXT([]byte(villages_depends))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(di.Depends))
	assert.Equal(t, "handle_schematics", di.Depends[0])
	assert.Equal(t, "default", di.Depends[1])
}
