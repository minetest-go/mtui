package web

import (
	"mtui/types"
	"net/http"
)

type MediaStats struct {
	Size             int64 `json:"size"`
	Count            int   `json:"count"`
	TransferredBytes int64 `json:"transferredbytes"`
}

func (a *Api) GetMediaStats(w http.ResponseWriter, r *http.Request) {
	SendJson(w, &MediaStats{
		Size:             a.app.Mediaserver.Size,
		Count:            a.app.Mediaserver.Count,
		TransferredBytes: a.app.Mediaserver.TransferredBytes.Load(),
	})
}

func (a *Api) ScanMedia(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	Send(w, true, a.app.Mediaserver.ScanDefaultSubdirs(a.app.WorldDir))
}
