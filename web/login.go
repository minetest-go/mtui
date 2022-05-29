package web

import (
	"encoding/json"
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

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	})

	token, err := t.SignedString([]byte("mykey"))

	res := &LoginResponse{
		Token: token,
	}

	Send(w, res, err)
}
