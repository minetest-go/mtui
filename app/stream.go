package app

import (
	"archive/zip"
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
