package app

import (
	"archive/zip"
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
	"path"
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

	zw := zip.NewWriter(w)
	defer zw.Close()

	bytes := int64(0)
	files := int64(0)
	buf := make([]byte, 1024*1024*10) // 10 mb buffer

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

		if IsSqliteDatabase(filePath) && !a.MaintenanceMode() {
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
		defer fsFile.Close()

		fc, err := io.CopyBuffer(zipFile, fsFile, buf)
		if err != nil {
			return err
		}
		bytes += fc
		files += 1

		return nil
	})

	return bytes, err
}

type DownloadZipOpts struct {
	Callback StreamProgressCallback
}

func (a *App) GetUncompressedZipSize(filename string) (int64, error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return 0, fmt.Errorf("openreader errror: %v", err)
	}
	defer zr.Close()

	// total uncompressed size
	uncompressed_bytes := uint64(0)
	for _, f := range zr.File {
		uncompressed_bytes += f.UncompressedSize64
	}
	return int64(uncompressed_bytes), nil
}

func (a *App) Unzip(abspath string, filename string, req *http.Request, c *types.Claims, opts *DownloadZipOpts) (int64, error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return 0, fmt.Errorf("openreader errror: %v", err)
	}
	defer zr.Close()

	bytes := int64(0)
	files := int64(0)

	for _, f := range zr.File {
		targetfile := path.Join(abspath, f.Name)
		dirname := path.Dir(targetfile)
		err = os.MkdirAll(dirname, 0644)
		if err != nil {
			return 0, fmt.Errorf("mkdirall error: %v", err)
		}

		if f.FileInfo().IsDir() {
			continue
		}

		zipfile, err := f.Open()
		if err != nil {
			return 0, fmt.Errorf("file open error: '%s', '%v'", f.Name, err)
		}

		if opts.Callback != nil {
			opts.Callback(files, bytes, f.Name)
		}

		fc, err := a.WriteFile(targetfile, zipfile, req, c)
		bytes += fc
		files++
		zipfile.Close()

		if err != nil {
			return 0, fmt.Errorf("writefile error: %v", err)
		}
	}

	return bytes, nil
}

func (a *App) DownloadToTempfile(r io.Reader) (string, error) {
	tf, err := os.CreateTemp(os.TempDir(), "mtui-upload")
	if err != nil {
		return "", fmt.Errorf("createtemp error: %v", err)
	}

	buf := make([]byte, 1024*1024*1) // 1 mb buffer
	_, err = io.CopyBuffer(tf, r, buf)
	if err != nil {
		return "", fmt.Errorf("copybuffer error: %v", err)
	}

	return tf.Name(), nil
}

func (a *App) DownloadAndUnzip(abspath string, r io.Reader, req *http.Request, c *types.Claims, opts *DownloadZipOpts) (int64, error) {

	if opts == nil {
		opts = &DownloadZipOpts{}
	}

	tempfile, err := a.DownloadToTempfile(r)
	if err != nil {
		return 0, fmt.Errorf("temp file error: %v", err)
	}

	return a.Unzip(abspath, tempfile, req, c, opts)
}
