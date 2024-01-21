package web

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "api/proxy") {
			// don't handle websocket connections
			next.ServeHTTP(w, r)
			return
		}

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		origin := r.Header.Get("X-Forwarded-For")
		if origin == "" {
			origin = r.RemoteAddr
		}

		logrus.WithFields(logrus.Fields{
			"origin":         origin,
			"method":         r.Method,
			"path":           r.URL.Path,
			"status":         rw.statusCode,
			"content-length": r.ContentLength,
			"user-agent":     r.Header.Get("User-Agent"),
		}).Debug("HTTP Request")
	})
}
