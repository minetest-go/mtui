package mail_test

import (
	"mtui/mail"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListMailNewDir(t *testing.T) {
	td, err := os.MkdirTemp(os.TempDir(), "mail")
	assert.NoError(t, err)
	assert.NoError(t, os.MkdirAll(path.Join(td, "mails", "contacts"), 0755))

	m := mail.New(td)
	assert.NotNil(t, m)

	list := make([]*mail.Message, 1)
	list[0] = &mail.Message{
		Body:    "abc",
		Sender:  "def",
		Subject: "subj",
		Time:    123,
		Unread:  true,
	}
	assert.NoError(t, m.SetMessages("singleplayer", list))

	list2, err := m.GetMessages("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, list2)
	assert.Equal(t, 1, len(list2))
	assert.Equal(t, list2[0].Body, list[0].Body)
	assert.Equal(t, list2[0].Sender, list[0].Sender)
	assert.Equal(t, list2[0].Subject, list[0].Subject)
	assert.Equal(t, list2[0].Time, list[0].Time)
	assert.Equal(t, list2[0].Unread, list[0].Unread)
}

func TestListMail(t *testing.T) {
	m := mail.New("testdata")
	assert.NotNil(t, m)

	list, err := m.GetMessages("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "mail\nbody", list[0].Body)
	assert.Equal(t, "otherplayer", list[0].Sender)
	assert.Equal(t, "subj", list[0].Subject)
	assert.Equal(t, int64(1563301643), list[0].Time)
	assert.Equal(t, false, list[0].Unread)

	list, err = m.GetMessages("unknownuser")
	assert.NoError(t, err)
	assert.Nil(t, list)
}

func TestContacts(t *testing.T) {
	m := mail.New("testdata")
	assert.NotNil(t, m)

	c, err := m.GetContacts("singleplayer")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c["admin"])
	assert.Equal(t, "admin", c["admin"].Name)
	assert.Equal(t, "", c["admin"].Note)

	c, err = m.GetContacts("unknownuser")
	assert.NoError(t, err)
	assert.Nil(t, c)
}
