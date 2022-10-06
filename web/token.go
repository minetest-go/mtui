package web

import (
	"errors"
	"mtui/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TOKEN_COOKIE_NAME = "mtui"

var err_unauthorized = errors.New("unauthorized")

func (api *Api) createCookie(value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    value,
		Path:     api.app.Config.CookiePath,
		Expires:  expires,
		Domain:   api.app.Config.CookieDomain,
		HttpOnly: true,
		Secure:   api.app.Config.CookieSecure,
	}
}

func (api *Api) SetToken(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, api.createCookie(token, expires))
}

func GetToken(r *http.Request) (string, error) {
	c, err := r.Cookie(TOKEN_COOKIE_NAME)
	if err == http.ErrNoCookie {
		return "", err_unauthorized
	}
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

func (api *Api) RemoveClaims(w http.ResponseWriter) {
	http.SetCookie(w, api.createCookie("", time.Now()))
}

func (api *Api) SetClaims(w http.ResponseWriter, claims *types.Claims) error {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(api.app.Config.JWTKey))
	if err != nil {
		return err
	}

	api.SetToken(w, token, claims.ExpiresAt.Local())
	return nil
}

func (api *Api) GetClaims(r *http.Request) (*types.Claims, error) {
	t, err := GetToken(r)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, err_unauthorized
	}

	token, err := jwt.ParseWithClaims(t, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(api.app.Config.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err_unauthorized
	}

	claims, ok := token.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	return claims, nil
}
