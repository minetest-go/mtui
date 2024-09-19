package web

import (
	"archive/zip"
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
	"path"
)

func (a *Api) UploadZip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	tf, err := os.CreateTemp(os.TempDir(), "mtui-zip-upload")
	if err != nil {
		SendError(w, 500, err)
		return
	}
	defer os.Remove(tf.Name())

	buf := make([]byte, 1024*1024*1) // 1 mb buffer
	_, err = io.CopyBuffer(tf, r.Body, buf)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	zr, err := zip.OpenReader(tf.Name())
	if err != nil {
		SendError(w, 500, err)
		return
	}
	defer zr.Close()

	count := int64(0)
	for _, f := range zr.File {
		targetfile := path.Join(absdir, f.Name)
		dirname := path.Dir(targetfile)
		err = os.MkdirAll(dirname, 0644)
		if err != nil {
			SendError(w, 500, err)
			return
		}

		if f.FileInfo().IsDir() {
			continue
		}

		zipfile, err := f.Open()
		if err != nil {
			SendError(w, 500, err)
			return
		}

		fc, err := a.app.WriteFile(targetfile, zipfile, r, claims)
		count += fc
		zipfile.Close()

		if err != nil {
			SendError(w, 500, err)
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
		SendError(w, 500, err)
		return
	}

	count, err := a.app.DownloadTargGZ(absdir, r.Body, r, claims, nil)
	if err != nil {
		SendError(w, 500, err)
		return
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
