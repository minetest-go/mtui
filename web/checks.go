package web

import (
	"mtui/types"
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
		} else if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		fn(w, r, claims)
	}
}

func SecurePriv(required_priv string, fn SecureHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaims(r)
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
