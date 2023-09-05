package web

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

type BrowseItem struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime int64  `json:"mtime"`
}

type BrowseResult struct {
	Dir       string        `json:"dir"`
	ParentDir string        `json:"parent_dir"`
	Items     []*BrowseItem `json:"items"`
}

func (a *Api) get_sanitized_dir(r *http.Request) (string, string, error) {
	dir := r.URL.Query().Get("dir")
	if strings.Contains(dir, "..") {
		return "", "", fmt.Errorf("invalid dir: '%s'", dir)
	}

	return dir, path.Join(a.app.WorldDir, dir), nil
}

func (a *Api) get_sanitized_filename(r *http.Request, query_param string) (string, error) {
	filename := r.URL.Query().Get(query_param)
	if strings.Contains(filename, "..") {
		return "", fmt.Errorf("invalid filename: '%s'", filename)
	}

	return path.Join(a.app.WorldDir, filename), nil
}

func (a *Api) BrowseFolder(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	entries, err := os.ReadDir(absdir)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	result := &BrowseResult{
		Dir:   reldir,
		Items: []*BrowseItem{},
	}

	if reldir != "/" {
		result.ParentDir = path.Dir(reldir)
	}

	for _, entry := range entries {

		if entry.IsDir() {
			result.Items = append(result.Items, &BrowseItem{
				Name:  entry.Name(),
				IsDir: true,
			})
		} else {
			item := &BrowseItem{
				Name:  entry.Name(),
				IsDir: false,
			}

			fi, _ := os.Stat(path.Join(absdir, entry.Name()))
			if fi != nil {
				item.Size = fi.Size()
				item.ModTime = fi.ModTime().Unix()
			}

			result.Items = append(result.Items, item)
		}
	}

	Send(w, result, nil)
}

func (a *Api) DownloadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	filename, err := a.get_sanitized_filename(r, "filename")
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
	_, err = f.Read(header)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	contentType := http.DetectContentType(header)
	w.Header().Set("Content-Type", contentType)

	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(filename)))
	}

	_, err = w.Write(header)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	_, err = io.Copy(w, f)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\".zip", zipfilename))
	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

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
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (a *Api) UploadZip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	_, absdir, err := a.get_sanitized_dir(r)
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

	for _, f := range zr.File {
		targetfile := path.Join(absdir, f.Name)
		dirname := path.Dir(targetfile)
		err = os.MkdirAll(dirname, 0644)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		targetf, err := os.OpenFile(targetfile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		zipfile, err := f.Open()
		if err != nil {
			targetf.Close()
			SendError(w, 500, err.Error())
			return
		}

		_, err = io.Copy(targetf, zipfile)
		targetf.Close()
		zipfile.Close()

		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}
}

func (a *Api) Mkdir(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	_, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.MkdirAll(absdir, 0644)
	Send(w, true, nil)
}

func (a *Api) UploadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	filename, err := a.get_sanitized_filename(r, "filename")
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

	_, err = io.Copy(f, r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (a *Api) DeleteFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.RemoveAll(filename)
	Send(w, true, err)
}

func (a *Api) RenameFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	src, err := a.get_sanitized_filename(r, "src")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	dst, err := a.get_sanitized_filename(r, "dst")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.Rename(src, dst)
	Send(w, true, err)
}
