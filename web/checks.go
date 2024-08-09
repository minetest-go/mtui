package web

import (
	"fmt"
	"mtui/types"
	"net/http"
	"sync/atomic"
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
			SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		} else if err != nil {
			api.RemoveClaims(w)
			SendError(w, http.StatusInternalServerError, err)
			return
		}
		fn(w, r, claims)
	}
}

// Optional login
func (api *Api) OptionalSecure(fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			// not logged in
			fn(w, r, nil)
			return
		} else if err != nil {
			api.RemoveClaims(w)
			SendError(w, http.StatusInternalServerError, err)
			return
		}
		// logged in
		fn(w, r, claims)
	}
}

type Check func(w http.ResponseWriter, r *http.Request) bool

func (api *Api) PrivCheck(required_priv string) Check {
	return func(w http.ResponseWriter, r *http.Request) bool {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return false
		} else if err != nil {
			api.RemoveClaims(w)
			SendError(w, http.StatusInternalServerError, err)
			return false
		}
		if !claims.HasPriv(required_priv) {
			SendError(w, http.StatusForbidden, fmt.Errorf("forbidden, missing priv: '%s'", required_priv))
			return false
		}
		return true
	}
}

func (api *Api) FeatureCheck(name string) Check {
	return func(w http.ResponseWriter, r *http.Request) bool {
		feature, err := api.app.Repos.FeatureRepository.GetByName(name)
		if err != nil {
			SendError(w, 500, err)
			return false
		}

		if !feature.Enabled {
			SendError(w, http.StatusInternalServerError, fmt.Errorf("Feature '%s' not enabled", name))
			return false
		}

		return true
	}
}

// maintenance mode

type MaintenanceModeCheckHandler struct {
	maint_mode *atomic.Bool
	handler    http.Handler
}

func (h MaintenanceModeCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.maint_mode.Load() {
		h.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Maintenance mode active"))
	}
}

func MaintenanceModeCheck(maint_mode *atomic.Bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return MaintenanceModeCheckHandler{maint_mode: maint_mode, handler: h}
	}
}

type SecureHandlerImpl struct {
	checks  []Check
	handler http.Handler
}

// secure handler with checks

func (sh SecureHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, check := range sh.checks {
		success := check(w, r)
		if !success {
			return
		}
	}
	sh.handler.ServeHTTP(w, r)
}

func SecureHandler(checks ...Check) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return SecureHandlerImpl{checks: checks, handler: h}
	}
}

// check for priv (UI)
func (api *Api) SecurePriv(required_priv string, fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := api.GetClaims(r)
		if err == err_unauthorized {
			SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		} else if err != nil {
			SendError(w, http.StatusInternalServerError, err)
			return
		}
		if !claims.HasPriv(required_priv) {
			SendError(w, http.StatusForbidden, fmt.Errorf("forbidden, missing priv: '%s'", required_priv))
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
			SendError(w, 500, err)
			return
		}

		if feature.Enabled {
			fn(w, r)
		} else {
			SendError(w, http.StatusInternalServerError, fmt.Errorf("Feature '%s' not enabled", name))
		}
	}
}
