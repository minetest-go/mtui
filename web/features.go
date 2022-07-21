package web

import (
	"net/http"
	"os"
	"path"
)

type Features struct {
	Mail    bool `json:"mail"`
	Areas   bool `json:"areas"`
	SkinsDB bool `json:"skinsdb"`
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

func checkFile(path string) bool {
	f, err := exists(path)
	if err != nil {
		// don't set flag if an error occurs
		return false
	}
	return f
}

func (a *Api) GetFeatures(w http.ResponseWriter, r *http.Request) {
	SendJson(w, &Features{
		Mail:    checkFile(path.Join(a.app.WorldDir, "mails")),
		Areas:   checkFile(path.Join(a.app.WorldDir, "areas.json")),
		SkinsDB: checkFile(path.Join(a.app.WorldDir, "worldmods", "skinsdb")),
	})
}
