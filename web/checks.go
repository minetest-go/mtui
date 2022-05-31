package web

import (
	"mtadmin/types"
	"net/http"
)

func CheckApiKey(key string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") != key {
			// unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// delegate
		fn(w, r)
	}
}

type SecureHandlerFunc func(http.ResponseWriter, *http.Request, *types.Claims)

func Secure(fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaims(r)
		if err == err_unauthorized {
			SendError(w, 401, "unauthorized")
			return
		}
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		fn(w, r, claims)
	}
}
