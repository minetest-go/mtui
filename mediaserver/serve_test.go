package mediaserver_test

import (
	"bytes"
	"encoding/hex"
	"mtui/mediaserver"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	m := mediaserver.New()
	wd, err := os.Getwd()
	dir := path.Join(wd, "testdata")
	assert.NoError(t, err)
	m.Scan(dir)

	hash_str := "7110eda4d09e062aa5e4a390b0a572ac0d2c0220"
	hash, err := hex.DecodeString(hash_str)
	assert.NoError(t, err)

	buf := append([]byte("MTHS"), 0, 1)
	buf = append(buf, hash...)

	// POST

	r := httptest.NewRequest(http.MethodPost, "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()

	m.ServeHTTPIndex(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, 26, w.Body.Len())
	body := w.Body.Bytes()

	assert.Equal(t, []byte("MTHS"), body[0:4])
	assert.Equal(t, uint8(0), body[4])
	assert.Equal(t, uint8(1), body[5])

	hash2 := hex.EncodeToString(body[6:])
	assert.Equal(t, hash_str, hash2)

	// GET

	r = httptest.NewRequest(http.MethodGet, "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"hash": hash_str,
	})

	m.ServeHTTPFetch(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, 4, w.Body.Len())
}
