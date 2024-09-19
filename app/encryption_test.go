package app_test

import (
	"bytes"
	"io"
	"mtui/app"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	b := []byte("hello world")
	key := "mykey"

	// write encrypted
	encrypted := bytes.NewBuffer([]byte{})
	w, err := app.EncryptedWriter(key, encrypted)
	assert.NoError(t, err)
	_, err = io.Copy(w, bytes.NewReader(b))
	assert.NoError(t, err)

	// read encrypted
	r, err := app.EncryptedReader(key, encrypted)
	assert.NoError(t, err)
	b2 := bytes.NewBuffer([]byte{})
	_, err = io.Copy(b2, r)
	assert.NoError(t, err)

	assert.Equal(t, b, b2.Bytes())
}
