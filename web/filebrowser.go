package web

import (
	"fmt"
	"mtui/app"
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
		SendError(w, 500, err)
		return
	}

	entries, err := os.ReadDir(absdir)
	if err != nil {
		SendError(w, 500, err)
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

func (a *Api) Mkdir(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	reldir, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	err = os.MkdirAll(absdir, 0644)
	Send(w, true, err)

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' created the directory '%s'", claims.Username, reldir),
	}, r)
}

func (a *Api) DeleteFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	err = os.RemoveAll(filename)
	Send(w, true, err)

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' deleted the file '%s'", claims.Username, rel_filename),
	}, r)
}

func (a *Api) UnzipFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_filename, filename, err := a.get_sanitized_filename(r, "filename")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		SendError(w, 500, fmt.Errorf("os open error: %v", err))
		return
	}
	defer f.Close()

	abspath := path.Dir(filename)
	count, err := a.app.DownloadZip(abspath, f, r, claims)
	Send(w, true, err)

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' unzipped the file '%s' (%d bytes)", claims.Username, rel_filename, count),
	}, r)
}

func (a *Api) RenameFile(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rel_src, src, err := a.get_sanitized_filename(r, "src")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	rel_dst, dst, err := a.get_sanitized_filename(r, "dst")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if app.IsDatabaseFile(dst) {
		// go into maintenance-mode during rename
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "maintenance",
			Message:  fmt.Sprintf("User '%s' renames '%s' to '%s' enabling temporary maintenance mode", claims.Username, rel_src, rel_dst),
		}, r)

		a.app.MaintenanceMode.Store(true)
		a.app.DetachDatabase()

		err = os.Rename(src, dst)

		a.app.AttachDatabase()
		a.app.MaintenanceMode.Store(false)
	} else {
		// no maintenance mode needed
		err = os.Rename(src, dst)
	}

	Send(w, true, err)

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "filebrowser",
		Message:  fmt.Sprintf("User '%s' moved the file '%s' to '%s'", claims.Username, rel_src, rel_dst),
	}, r)
}

func (a *Api) GetDirectorySize(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	_, absdir, err := a.get_sanitized_dir(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count := int64(0)
	err = filepath.Walk(absdir, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if app.IgnoreSqliteFileDownload(filePath) {
			return nil
		}

		fi, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("stat error: '%s': %v", filePath, err)
		}

		count += fi.Size()
		return nil
	})

	Send(w, count, err)
}
