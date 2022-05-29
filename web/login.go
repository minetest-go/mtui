package web

import (
	"encoding/json"
	"mtadmin/auth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	*jwt.StandardClaims
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	auth_entry, err := a.app.Repos.Auth.GetByUsername(req.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth_entry == nil {
		SendError(w, 404, "user not found")
		return
	}

	salt, verifier, err := auth.ParseDBPassword(auth_entry.Password)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	ok, err := auth.VerifyAuth(req.Username, req.Password, salt, verifier)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if !ok {
		SendError(w, 401, "unauthorized")
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	})

	token, err := t.SignedString([]byte("mykey"))

	http.SetCookie(w, &http.Cookie{
		Name:     "mtadmin",
		Value:    token,
		Path:     "/", //TODO: configure
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Domain:   "127.0.0.1", //TODO
		HttpOnly: true,
		Secure:   true, //TODO
	})

	w.WriteHeader(http.StatusOK)
}
