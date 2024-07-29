package web

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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

func (a *Api) writeFile(filename string, data io.Reader, r *http.Request, claims *types.Claims) (int64, error) {

	is_dbfile := usedDatabaseFiles[path.Base(filename)]
	if is_dbfile && !a.app.MaintenanceMode.Load() {
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "maintenance",
			Message:  fmt.Sprintf("User '%s' uploads '%s' enabling temporary maintenance mode", claims.Username, path.Base(filename)),
		}, r)

		// not in maintenance mode, enable it for the upload of this critical database file
		a.app.MaintenanceMode.Store(true)
		a.app.DetachDatabase()
		defer func() {
			a.app.AttachDatabase()
			a.app.MaintenanceMode.Store(false)
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

func (a *Api) UploadZip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	tf, err := os.CreateTemp(os.TempDir(), "mtui-zip-upload")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	defer os.Remove(tf.Name())

	_, err = io.Copy(tf, r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	zr, err := zip.OpenReader(tf.Name())
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	defer zr.Close()

	count := int64(0)
	for _, f := range zr.File {
		targetfile := path.Join(absdir, f.Name)
		dirname := path.Dir(targetfile)
		err = os.MkdirAll(dirname, 0644)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		if f.FileInfo().IsDir() {
			continue
		}

		zipfile, err := f.Open()
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		fc, err := a.writeFile(targetfile, zipfile, r, claims)
		count += fc
		zipfile.Close()

		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded a zip to the directory '%s' with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) UploadTarGZ(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	zr, err := gzip.NewReader(r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	defer zr.Close()

	tr := tar.NewReader(zr)
	count := int64(0)

	for {
		th, err := tr.Next()
		if th == nil || err == io.EOF {
			break
		}
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		targetfile := path.Join(absdir, th.Name)
		dirname := path.Dir(targetfile)
		err = os.MkdirAll(dirname, 0644)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		if th.FileInfo().IsDir() {
			continue
		}

		fc, err := a.writeFile(targetfile, tr, r, claims)
		count += fc

		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded a tar.gz to the directory '%s' with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) UploadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	count, err := a.writeFile(filename, r.Body, r, claims)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded the file '%s' with %d bytes", claims.Username, rel_filename, count),
	}, r)
}
