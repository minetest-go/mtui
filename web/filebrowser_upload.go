package web

import (
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
)

func (a *Api) UploadZip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := a.app.DownloadZip(absdir, r.Body, r, claims, nil)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded a zip to the directory '%s' with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) UploadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := a.app.WriteFile(filename, r.Body, r, claims)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded the file '%s' with %d bytes", claims.Username, rel_filename, count),
	}, r)
}

func (a *Api) AppendFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	_, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		SendError(w, 500, fmt.Errorf("openfile error for '%s': %v", filename, err))
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		SendError(w, 500, fmt.Errorf("copyfile error for '%s': %v", filename, err))
		return
	}
}
