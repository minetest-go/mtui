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
	"path/filepath"
	"strings"
)

func (a *Api) DownloadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	defer f.Close()

	header := make([]byte, 512)
	header_size, err := f.Read(header)
	if err != nil && err != io.EOF {
		SendError(w, 500, err.Error())
		return
	}
	contentType := http.DetectContentType(header)
	w.Header().Set("Content-Type", contentType)

	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(filename)))
	}

	_, err = w.Write(header[:header_size])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	count, err := io.Copy(w, f)
	if err != nil {
		SendError(w, 500, err.Error())
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
		SendError(w, 500, err.Error())
		return
	}

	zipfilename := path.Base(absdir)
	if reldir == "/" {
		zipfilename = "world"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", zipfilename))
	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	count := int64(0)
	err = filepath.Walk(absdir, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, absdir)
		zipFile, err := zw.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		fc, err := io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		count += fc
		return nil
	})

	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' as zip with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) DownloadTarGZ(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	targzfilename := path.Base(absdir)
	if reldir == "/" {
		targzfilename = "world"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar.gz\"", targzfilename))
	w.Header().Set("Content-Type", "application/gzip")

	zw := gzip.NewWriter(w)
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	count := int64(0)
	err = filepath.Walk(absdir, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, absdir)
		fi, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		fi.Name = relPath

		err = tw.WriteHeader(fi)
		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		fc, err := io.Copy(tw, fsFile)
		if err != nil {
			return err
		}
		count += fc
		return nil
	})

	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' as tar.gz with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}
