package web

import "net/http"

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
