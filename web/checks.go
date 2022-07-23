package web

import (
	"mtui/types"
	"net/http"
)

func (api *Api) CheckApiKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") != api.app.Config.APIKey {
			// unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// delegate
		fn(w, r)
	}
}

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)

func (api *Api) Secure(fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, 401, "unauthorized")
			return
		} else if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		fn(w, r, claims)
	}
}

func (api *Api) SecurePriv(required_priv string, fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, 401, "unauthorized")
			return
		} else if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		has_priv := false
		for _, priv := range claims.Privileges {
			if priv == required_priv {
				has_priv = true
			}
		}
		if !has_priv {
			SendError(w, 403, "forbidden, missing priv: "+required_priv)
			return
		}
		fn(w, r, claims)
	}
}
