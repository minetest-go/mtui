package web

import (
	"mtui/public"
	"mtui/types"
	"net/http"
)

func (a *Api) GetThemes(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	themes := []string{}
	for name := range public.ThemeMap {
		themes = append(themes, name)
	}
	Send(w, themes, nil)
}
