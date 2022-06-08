package mail_test

import (
	"mtui/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListMail(t *testing.T) {
	m := mail.New("testdata")
	assert.NotNil(t, m)
}
