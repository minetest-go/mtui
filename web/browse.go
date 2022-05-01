package web

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"mtadmin/app"
	"net/http"

	"github.com/gorilla/mux"
)

type Browse struct {
	App *app.App
}

type BrowseRequest struct {
	Path string `json:"path"`
}

type BrowseResponse struct {
	Path        string         `json:"path"`
	Directories []string       `json:"directories"`
	Files       []FileResponse `json:"files"`
}

type FileResponse struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	CreateTime int64  `json:"ctime"`
	ApiPath    string `json:"apipath"`
}

func SetupBrowse(r *mux.Router, a *app.App) *Browse {
	b := &Browse{App: a}
	r.HandleFunc("/api/browse", b.HandleBrowse).Methods(http.MethodPost)
	return b
}

func scanDir(worldpath, path string, res *BrowseResponse) error {
	list, err := ioutil.ReadDir(worldpath + path)
	if err != nil {
		return err
	}

	for _, entry := range list {
		if entry.IsDir() {
			res.Directories = append(res.Directories, entry.Name())

		} else if entry.Mode().IsRegular() {
			encoded_path := base64.StdEncoding.EncodeToString([]byte(path + "/" + entry.Name()))

			res.Files = append(res.Files, FileResponse{
				Name:       entry.Name(),
				Path:       path,
				Size:       entry.Size(),
				CreateTime: entry.ModTime().Unix(),
				ApiPath:    "api/browse/" + encoded_path,
			})

		}
	}

	return nil
}

func (b *Browse) HandleBrowse(w http.ResponseWriter, r *http.Request) {
	browse_req := &BrowseRequest{}
	err := json.NewDecoder(r.Body).Decode(browse_req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	res := &BrowseResponse{
		Path:        browse_req.Path,
		Directories: []string{},
		Files:       []FileResponse{},
	}

	err = scanDir(b.App.WorldDir, browse_req.Path, res)
	Send(w, res, err)
}
