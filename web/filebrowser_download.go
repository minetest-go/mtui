package web

import (
	"fmt"
	"io"
	"mtui/app"
	"mtui/types"
	"net/http"
	"os"
	"path"
)

func (a *Api) DownloadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if app.IsSqliteDatabase(filename) && !a.app.MaintenanceMode() {
		tmppath, err := app.CreateSqliteSnapshot(filename)
		if err != nil {
			SendError(w, 500, fmt.Errorf("error creating snapshot of '%s': %v", filename, err))
			return
		}
		defer os.Remove(tmppath)
		filename = tmppath
	}

	f, err := os.Open(filename)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	defer f.Close()

	header := make([]byte, 512)
	header_size, err := f.Read(header)
	if err != nil && err != io.EOF {
		SendError(w, 500, err)
		return
	}
	contentType := http.DetectContentType(header)
	w.Header().Set("Content-Type", contentType)

	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(rel_filename)))
	}

	_, err = w.Write(header[:header_size])
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := io.Copy(w, f)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the file '%s' with %d bytes", claims.Username, rel_filename, count+int64(header_size)),
	}, r)
}

func (a *Api) DownloadZip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	zipfilename := path.Base(absdir)
	if reldir == "/" {
		zipfilename = "world"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", zipfilename))
	w.Header().Set("Content-Type", "application/zip")

	count, err := a.app.StreamZip(absdir, w, nil)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' as zip with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}
