package mail_test

import (
	"mtui/mail"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/minetest-go/mtdb"
	"github.com/stretchr/testify/assert"
)

func TestListMailNewDir(t *testing.T) {
	td, err := os.MkdirTemp(os.TempDir(), "mail")
	assert.NoError(t, err)

	contents := `
backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
mod_storage_backend = sqlite3
`
	err = os.WriteFile(path.Join(td, "world.mt"), []byte(contents), 0644)
	assert.NoError(t, err)

	ctx, err := mtdb.New(td)
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	m := mail.New(ctx)
	assert.NotNil(t, m)

	// get empty
	entry, err := m.GetEntry("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, 0, len(entry.Contacts))
	assert.Equal(t, 0, len(entry.Drafts))
	assert.Equal(t, 0, len(entry.Inbox))
	assert.Equal(t, 0, len(entry.Outbox))
	assert.Equal(t, 0, len(entry.Lists))

	// add mail
	entry.Inbox = append(entry.Inbox, &mail.Message{
		ID:      uuid.NewString(),
		From:    "xy",
		To:      "singleplayer",
		Subject: "abc",
		Body:    "body",
	})

	// set
	assert.NoError(t, m.SetEntry("singleplayer", entry))

	// get updated entry
	entry, err = m.GetEntry("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, 1, len(entry.Inbox))
	assert.Equal(t, "xy", entry.Inbox[0].From)
	assert.Equal(t, "singleplayer", entry.Inbox[0].To)
	assert.Equal(t, "abc", entry.Inbox[0].Subject)
	assert.Equal(t, "body", entry.Inbox[0].Body)

}
