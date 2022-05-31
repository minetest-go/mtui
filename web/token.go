package web

import (
	"errors"
	"mtadmin/types"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TOKEN_COOKIE_NAME = "mtadmin"

var err_unauthorized = errors.New("unauthorized")

func createCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    value,
		Path:     os.Getenv("COOKIE_PATH"),
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		HttpOnly: true,
		Secure:   os.Getenv("COOKIE_SECURE") == "true",
	}
}

func SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, createCookie(token))
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

func RemoveClaims(w http.ResponseWriter) {
	http.SetCookie(w, createCookie(""))
}

func SetClaims(w http.ResponseWriter, claims *types.Claims) error {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		return err
	}

	SetToken(w, token)
	return nil
}

func GetClaims(r *http.Request) (*types.Claims, error) {
	t, err := GetToken(r)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, err_unauthorized
	}

	token, err := jwt.ParseWithClaims(t, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTKEY")), nil
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
