package web

import (
	"bytes"
	"encoding/json"
	"mtadmin/app"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestBrowse(t *testing.T) {
	// setup env
	worlddir, err := os.MkdirTemp(os.TempDir(), "world")
	assert.NoError(t, err)
	assert.NotNil(t, worlddir)

	router := mux.NewRouter()
	a, err := app.Create(worlddir)
	assert.NoError(t, err)

	browse := SetupBrowse(router, a)
	assert.NotNil(t, browse)

	// setup request
	data, err := json.Marshal(&BrowseRequest{
		Path: "/",
	})
	assert.NoError(t, err)
	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	// request
	browse.HandleBrowse(w, r)

	// verify
	assert.Equal(t, 200, w.Result().StatusCode)

	result := &BrowseResponse{}
	assert.NoError(t, json.NewDecoder(w.Result().Body).Decode(result))

	assert.Equal(t, "/", result.Path)
	assert.True(t, len(result.Files) > 0)
	assert.Equal(t, 0, len(result.Directories))
}
