package web

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"mtui/db"
	"mtui/types"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ignoreFileDownload(filename string) bool {
	if strings.HasSuffix(filename, ".sqlite-shm") || strings.HasSuffix(filename, ".sqlite-wal") {
		// sqlite wal and shared memory
		return true
	}
	return false
}

func isSqliteDatabase(filename string) bool {
	return strings.HasSuffix(filename, ".sqlite")
}

func createSqliteSnapshot(filename string) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "backup.sqlite")
	if err != nil {
		return "", fmt.Errorf("create temp error: %v", err)
	}

	err = db.BackupSqlite3Database(context.Background(), filename, f.Name())
	if err != nil {
		return "", fmt.Errorf("backup error from '%s' to '%s': %v", filename, f.Name(), err)
	}

	return f.Name(), nil
}

func (a *Api) DownloadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	maintenance := a.app.MaintenanceMode.Load()
	if isSqliteDatabase(filename) && !maintenance {
		tmppath, err := createSqliteSnapshot(filename)
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
	maintenance := a.app.MaintenanceMode.Load()

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
		if ignoreFileDownload(filePath) {
			return nil
		}
		relPath := strings.TrimPrefix(filePath, absdir)
		relPath = strings.TrimPrefix(relPath, "/")

		if isSqliteDatabase(filePath) && !maintenance {
			tmppath, err := createSqliteSnapshot(filePath)
			if err != nil {
				return fmt.Errorf("sqlite snapshot error for '%s': %v", filePath, err)
			}
			defer os.Remove(tmppath)
			filePath = tmppath
		}

		zipFile, err := zw.CreateHeader(&zip.FileHeader{
			Name:     relPath,
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		})
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
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' as zip with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) DownloadTarGZ(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	maintenance := a.app.MaintenanceMode.Load()

	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
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
		if ignoreFileDownload(filePath) {
			return nil
		}

		relPath := strings.TrimPrefix(filePath, absdir)
		fi, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		fi.Name = relPath

		if isSqliteDatabase(filePath) && !maintenance {
			tmppath, err := createSqliteSnapshot(filePath)
			if err != nil {
				return fmt.Errorf("sqlite snapshot error for '%s': %v", filePath, err)
			}
			defer os.Remove(tmppath)
			filePath = tmppath

			tmpfi, err := os.Stat(tmppath)
			if err != nil {
				return fmt.Errorf("stat error: '%s': %v", tmppath, err)
			}

			fi.Size = tmpfi.Size()
		}

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
		SendError(w, 500, err)
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' as tar.gz with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}
