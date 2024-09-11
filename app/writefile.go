package app

import (
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
	"path"
)

var usedDatabaseFiles = map[string]bool{
	"map.sqlite":         true,
	"auth.sqlite":        true,
	"mod_storage.sqlite": true,
	"players.sqlite":     true,
	"mtui.sqlite":        true,
}

func (a *App) WriteFile(filename string, data io.Reader, r *http.Request, claims *types.Claims) (int64, error) {

	is_dbfile := usedDatabaseFiles[path.Base(filename)]
	if is_dbfile && !a.MaintenanceMode.Load() {
		// create log entry
		a.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "maintenance",
			Message:  fmt.Sprintf("User '%s' uploads '%s' enabling temporary maintenance mode", claims.Username, path.Base(filename)),
		}, r)

		// not in maintenance mode, enable it for the upload of this critical database file
		a.MaintenanceMode.Store(true)
		a.DetachDatabase()
		defer func() {
			a.AttachDatabase()
			a.MaintenanceMode.Store(false)
		}()
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("openfile error for '%s': %v", filename, err)
	}
	defer f.Close()

	count, err := io.Copy(f, data)
	if err != nil {
		return 0, fmt.Errorf("copyfile error for '%s': %v", filename, err)
	}

	return count, nil
}
