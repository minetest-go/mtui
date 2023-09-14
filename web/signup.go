package web

import (
	"encoding/json"
	"fmt"
	"mtui/auth"
	"mtui/types"
	"net/http"

	"github.com/dchest/captcha"
	dbauth "github.com/minetest-go/mtdb/auth"
)

type SignupRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

func (a *Api) SignupCaptcha(w http.ResponseWriter, r *http.Request) {
	SendText(w, captcha.New())
}

func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {
	sr := &SignupRequest{}
	err := json.NewDecoder(r.Body).Decode(sr)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if !captcha.VerifyString(sr.CaptchaID, sr.Captcha) {
		SendError(w, 401, "captcha invalid")
		return
	}

	err = auth.ValidateUsername(sr.Username)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("Username validation failed: %v", err))
		return
	}

	if sr.Password == "" {
		SendError(w, 500, "empty password")
		return
	}

	auth_entry, err := a.app.DBContext.Auth.GetByUsername(sr.Username)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("auth-db error: '%v'", err))
		return
	}

	if auth_entry != nil {
		SendError(w, 500, fmt.Sprintf("player already exists: '%s'", sr.Username))
		return
	}

	salt, verifier, err := auth.CreateAuth(sr.Username, sr.Password)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("create-auth failed: %v", err))
		return
	}

	dbstr := auth.CreateDBPassword(salt, verifier)
	auth_entry = &dbauth.AuthEntry{
		Name:     sr.Username,
		Password: dbstr,
	}
	err = a.app.DBContext.Auth.Create(auth_entry)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("authdb insert failed: %v", err))
		return
	}

	// TODO: configurable privs
	for _, priv := range []string{"interact", "shout"} {
		err = a.app.DBContext.Privs.Create(&dbauth.PrivilegeEntry{
			ID:        *auth_entry.ID,
			Privilege: priv,
		})

		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: sr.Username,
		Event:    "signup",
		Message:  fmt.Sprintf("User '%s' signed up successfully", sr.Username),
	}, r)
}
