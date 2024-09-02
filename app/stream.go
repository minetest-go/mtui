package app

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) StreamZip(path string, w io.Writer) (int64, error) {
	maintenance := a.MaintenanceMode.Load()

	zw := zip.NewWriter(w)
	defer zw.Close()

	count := int64(0)
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

	return count, err
}
