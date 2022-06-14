package web

import (
	"net/http"
	"os"
	"path"
)

type Features struct {
	Mail  bool `json:"mail"`
	Areas bool `json:"areas"`
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (a *Api) GetFeatures(w http.ResponseWriter, r *http.Request) {
	has_mail, err := exists(path.Join(a.app.WorldDir, "mails"))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	has_areas, err := exists(path.Join(a.app.WorldDir, "areas.json"))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, &Features{
		Mail:  has_mail,
		Areas: has_areas,
	})
}
