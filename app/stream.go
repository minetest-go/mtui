package app

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type StreamProgressCallback func(files, bytes int64, currentfile string)

type StreamZipOpts struct {
	Callback StreamProgressCallback
}

func (a *App) StreamZip(path string, w io.Writer, opts *StreamZipOpts) (int64, error) {
	if opts == nil {
		opts = &StreamZipOpts{}
	}

	maintenance := a.MaintenanceMode.Load()

	zw := zip.NewWriter(w)
	defer zw.Close()

	bytes := int64(0)
	files := int64(0)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if IgnoreSqliteFileDownload(filePath) {
			return nil
		}
		relPath := strings.TrimPrefix(filePath, path)
		relPath = strings.TrimPrefix(relPath, "/")

		if IsSqliteDatabase(filePath) && !maintenance {
			tmppath, err := CreateSqliteSnapshot(filePath)
			if err != nil {
				return fmt.Errorf("sqlite snapshot error for '%s': %v", filePath, err)
			}
			defer os.Remove(tmppath)
			filePath = tmppath
		}

		if opts.Callback != nil {
			opts.Callback(files, bytes, relPath)
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
		bytes += fc
		files += 1

		return nil
	})

	return bytes, err
}

type StreamTarGZOpts struct {
	Callback StreamProgressCallback
}

func (a *App) StreamTarGZ(path string, w io.Writer, opts *StreamTarGZOpts) (int64, error) {
	if opts == nil {
		opts = &StreamTarGZOpts{}
	}

	maintenance := a.MaintenanceMode.Load()

	zw := gzip.NewWriter(w)
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	bytes := int64(0)
	files := int64(0)

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if IgnoreSqliteFileDownload(filePath) {
			return nil
		}

		relPath := strings.TrimPrefix(filePath, path)
		fi, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		fi.Name = relPath

		if IsSqliteDatabase(filePath) && !maintenance {
			tmppath, err := CreateSqliteSnapshot(filePath)
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

		if opts.Callback != nil {
			opts.Callback(files, bytes, relPath)
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
		bytes += fc
		files += 1

		return nil
	})

	return bytes, err
}
