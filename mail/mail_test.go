package mail_test

import (
	"mtui/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListMail(t *testing.T) {
	m := mail.New("testdata")
	assert.NotNil(t, m)

	list, err := m.List("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "mail\nbody", list[0].Body)
	assert.Equal(t, "otherplayer", list[0].Sender)
	assert.Equal(t, "subj", list[0].Subject)
	assert.Equal(t, int64(1563301643), list[0].Time)
	assert.Equal(t, false, list[0].Unread)

	_, err = m.List("unknownuser")
	assert.Error(t, err)
}

func TestContacts(t *testing.T) {
	m := mail.New("testdata")
	assert.NotNil(t, m)

	c, err := m.Contacts("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c["admin"])
	assert.Equal(t, "admin", c["admin"].Name)
	assert.Equal(t, "", c["admin"].Note)

	_, err = m.Contacts("unknownuser")
	assert.Error(t, err)
}
