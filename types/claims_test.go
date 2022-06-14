package types_test

import (
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimsHasPriv(t *testing.T) {
	c := types.Claims{
		Privileges: []string{"a", "b"},
	}

	assert.True(t, c.HasPriv("a"))
	assert.True(t, c.HasPriv("b"))
	assert.False(t, c.HasPriv("c"))
}
