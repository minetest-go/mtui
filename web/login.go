package web

import (
	"encoding/json"
	"fmt"
	"mtui/auth"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var tan_map = make(map[string]string)

func (a *Api) TanSetListener(c chan *bridge.CommandResponse) {
	for {
		cmd := <-c
		tc := &command.TanCommand{}
		err := json.Unmarshal(cmd.Data, tc)
		if err != nil {
			fmt.Printf("Tan-listener-error: %s\n", err.Error())
			continue
		}

		if tc.TAN == "" {
			// remove tan
			delete(tan_map, tc.Playername)
		} else {
			// set tan
			tan_map[tc.Playername] = tc.TAN
		}
	}
}

func (a *Api) DoLogout(w http.ResponseWriter, r *http.Request) {
	a.RemoveClaims(w)
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetClaims(r)
	if err == err_unauthorized {
		SendError(w, 401, "unauthorized")
	} else if err != nil {
		SendError(w, 500, err.Error())
	} else {
		// refresh token
		auth_entry, err := a.app.DBContext.Auth.GetByUsername(claims.Username)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		if auth_entry == nil {
			SendError(w, 404, "auth entry not found")
			return
		}

		claims, err = a.updateToken(w, *auth_entry.ID, claims.Username)
		Send(w, claims, err)
	}
}

func (a *Api) updateToken(w http.ResponseWriter, id int64, username string) (*types.Claims, error) {
	privs, err := a.app.DBContext.Privs.GetByID(id)
	if err != nil {
		return nil, err
	}

	priv_arr := make([]string, len(privs))
	for i, p := range privs {
		priv_arr[i] = p.Privilege
	}

	expires := time.Now().Add(7 * 24 * time.Hour)
	claims := &types.Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		Username:   username,
		Privileges: priv_arr,
	}
	return claims, a.SetClaims(w, claims)
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

	tan := tan_map[req.Username]

	if tan == "" {
		// login against the database password
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
	} else {
		// login with tan
		if tan != req.Password {
			SendError(w, 401, "unauthorized")
			return
		}

		// remove tan (single-use)
		delete(tan_map, req.Username)
	}

	claims, err := a.updateToken(w, *auth_entry.ID, auth_entry.Name)
	Send(w, claims, err)
}
