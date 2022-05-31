package web

import (
	"encoding/json"
	"mtadmin/auth"
	"mtadmin/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaims(r)
	if err == err_unauthorized {
		SendError(w, 401, "unauthorized")
	} else {
		Send(w, claims, err)
	}
}

func (a *Api) DoLogin(w http.ResponseWriter, r *http.Request) {
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

	err = SetClaims(w, &types.Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
		Username: req.Username,
	})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
