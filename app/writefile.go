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

func IsDatabaseFile(filename string) bool {
	return usedDatabaseFiles[path.Base(filename)]
}

func (a *App) WriteFile(filename string, data io.Reader, r *http.Request, claims *types.Claims) (int64, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("openfile error for '%s': %v", filename, err)
	}
	defer f.Close()

	buf := make([]byte, 1024*1024) // 1 mb buffer
	count, err := io.CopyBuffer(f, data, buf)
	if err != nil {
		return 0, fmt.Errorf("copyfile error for '%s': %v", filename, err)
	}

	return count, nil
}
