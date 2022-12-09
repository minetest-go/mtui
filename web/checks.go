package web

import (
	"fmt"
	"mtui/types"
	"net/http"
)

// check api-key (for the minetest engine calls)
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

// check api-key or claims/jwt (for the engine _or_ the UI)
func (api *Api) CheckApiKeyOrPriv(required_priv string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") == api.app.Config.APIKey {
			fn(w, r)
			return
		}
		claims, err := api.GetClaims(r)
		if err == nil && claims != nil && claims.HasPriv(required_priv) {
			fn(w, r)
			return
		}
		// unauthorized
		w.WriteHeader(http.StatusUnauthorized)
	}
}

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)

// check for login only (UI)
func (api *Api) Secure(fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, "unauthorized")
			return
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fn(w, r, claims)
	}
}

// check for priv (UI)
func (api *Api) SecurePriv(required_priv string, fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, "unauthorized")
			return
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !claims.HasPriv(required_priv) {
			SendError(w, http.StatusForbidden, "forbidden, missing priv: "+required_priv)
			return
		}
		fn(w, r, claims)
	}
}

// check if a feature is enabled
func (api *Api) Feature(name string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feature, err := api.app.Repos.FeatureRepository.GetByName(name)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		if feature.Enabled {
			fn(w, r)
		} else {
			SendError(w, http.StatusInternalServerError, fmt.Sprintf("Feature '%s' not enabled", name))
		}
	}
}
