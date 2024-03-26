package web

import (
	"fmt"
	"io"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minetest-go/mtdb/mod_storage"
)

func (a *Api) GetMtUIStorage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	entry, err := a.app.DBContext.ModStorage.Get("mtui", []byte(key))
	if err != nil {
		SendError(w, 500, err.Error())
	} else if entry == nil {
		SendError(w, 404, "not found")
	} else {
		SendText(w, string(entry.Value))
	}
}

func (a *Api) SetMtUIStorage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("readall error: %s", err.Error()))
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "uimod_storage",
		Message:  fmt.Sprintf("User '%s' sets the ui storage entry '%s' to '%s'", claims.Username, key, string(value)),
	}, r)

	entry, err := a.app.DBContext.ModStorage.Get("mtui", []byte(key))
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get error: %s", err.Error()))
		return
	}

	if entry == nil {
		// create
		entry = &mod_storage.ModStorageEntry{
			ModName: "mtui",
			Key:     []byte(key),
			Value:   value,
		}
		err = a.app.DBContext.ModStorage.Create(entry)
	} else {
		// update
		entry.Value = value
		err = a.app.DBContext.ModStorage.Update(entry)
	}

	Send(w, entry.Value, err)
}
