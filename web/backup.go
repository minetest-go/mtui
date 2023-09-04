package web

import (
	"archive/zip"
	"mtui/types"
	"net/http"
)

func (a *Api) Export(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	z := zip.NewWriter(w)
	err := a.app.DBContext.Export(z)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
