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

func (a *Api) get_sanitized_filename(r *http.Request, query_param string) (string, string, error) {
	filename := r.URL.Query().Get(query_param)
	if strings.Contains(filename, "..") {
		return "", "", fmt.Errorf("invalid filename: '%s'", filename)
	}

	return filename, path.Join(a.app.WorldDir, filename), nil
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

	a.CreateUILogEntry(&types.Log{
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\".zip", zipfilename))
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

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' downloaded the directory '%s' with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
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

		fc, err := io.Copy(targetf, zipfile)
		count += fc
		targetf.Close()
		zipfile.Close()

		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded a zip to the directory '%s' with %d bytes (uncompressed)", claims.Username, reldir, count),
	}, r)
}

func (a *Api) Mkdir(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.MkdirAll(absdir, 0644)
	Send(w, true, err)

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' created the directory '%s'", claims.Username, reldir),
	}, r)
}

func (a *Api) UploadFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	defer f.Close()

	count, err := io.Copy(f, r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' uploaded the file '%s' with %d bytes", claims.Username, rel_filename, count),
	}, r)
}

func (a *Api) DeleteFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.RemoveAll(filename)
	Send(w, true, err)

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' deleted the file '%s'", claims.Username, rel_filename),
	}, r)
}

func (a *Api) RenameFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_src, src, err := a.get_sanitized_filename(r, "src")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	rel_dst, dst, err := a.get_sanitized_filename(r, "dst")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = os.Rename(src, dst)
	Send(w, true, err)

	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' moved the file '%s' to '%s'", claims.Username, rel_src, rel_dst),
	}, r)
}
