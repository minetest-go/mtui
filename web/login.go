package web

import (
	"encoding/json"
	"mtadmin/auth"
	"net/http"
	"os"
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
	*jwt.RegisteredClaims
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	auth_entry, err := a.app.DBContext.Auth.GetByUsername(req.Username)
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
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	})

	token, err := t.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SetToken(w, token)

	w.WriteHeader(http.StatusOK)
}
