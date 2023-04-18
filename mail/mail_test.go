package mail_test

import (
	"mtui/mail"
	"os"
	"testing"

	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestListMailNewDir(t *testing.T) {
	td, err := os.MkdirTemp(os.TempDir(), "mail")
	assert.NoError(t, err)

	ctx, err := mtdb.New(td)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	m := mail.New(ctx)
	assert.NotNil(t, m)

	entry, err := m.GetEntry("singleplayer")
	assert.NoError(t, err)
	assert.Nil(t, entry)

	assert.NoError(t, m.SetEntry("singleplayer", &mail.PlayerEntry{}))
}
