package web

import (
	"errors"
	"net/http"
	"time"
)

const TOKEN_COOKIE_NAME = "mtadmin"

var err_unauthorized = errors.New("unauthorized")

func SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "mtadmin",
		Value:    token,
		Path:     "/", //TODO: configure
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Domain:   "127.0.0.1", //TODO
		HttpOnly: true,
		Secure:   true, //TODO
	})
}

func GetToken(r *http.Request) (string, error) {
	c, err := r.Cookie(TOKEN_COOKIE_NAME)
	if err != nil {
		return "", err
	}

	if c == nil {
		return "", err_unauthorized
	}

	return c.Value, nil
}
