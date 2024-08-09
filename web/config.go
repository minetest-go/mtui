package web

import (
	"bytes"
	"fmt"
	"io"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetConfig(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	e, err := a.app.Repos.ConfigRepo.GetByKey(types.ConfigKey(key))
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if e != nil {
		w.Write([]byte(e.Value))
	}
}

func (a *Api) SetConfig(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	e, err := a.app.Repos.ConfigRepo.GetByKey(types.ConfigKey(key))
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if e == nil {
		e = &types.ConfigEntry{
			Key: types.ConfigKey(key),
		}
	}

	e.Value = buf.String()
	err = a.app.Repos.ConfigRepo.Set(e)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "lua",
		Message:  fmt.Sprintf("User '%s' sets the config '%s' to '%s'", c.Username, key, e.Value),
	}, r)

}
