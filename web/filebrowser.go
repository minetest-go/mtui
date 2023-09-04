package web

import (
	"fmt"
	"io"
	"mtui/types"
	"net/http"
	"os"
	"path"
	"strings"
)

type BrowseItem struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
}

type BrowseResult struct {
	Dir   string        `json:"dir"`
	Items []*BrowseItem `json:"items"`
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

	err = os.Remove(filename)
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
