package web

import (
	"encoding/json"
	"fmt"
	"mtui/auth"
	"mtui/types"
	"net/http"
	"sync/atomic"

	dbauth "github.com/minetest-go/mtdb/auth"
)

var onboardFinished atomic.Bool

func (a *Api) canOboard() (bool, error) {
	if onboardFinished.Load() {
		return false, nil
	}

	// count all users
	count, err := a.app.DBContext.Auth.Count(&dbauth.AuthSearch{})
	if err != nil {
		return false, err
	}

	if count == 0 {
		// no user in the database, onboarding enabled
		return true, nil
	}

	if count == 1 {
		// one user found, check if this is the singleplayer
		sp := "singleplayer"
		count, err = a.app.DBContext.Auth.Count(&dbauth.AuthSearch{Username: &sp})
		if err != nil {
			return false, err
		}

		if count == 1 {
			// only singleplayer found, allow onboarding
			return true, nil
		}
	}

	if count > 1 {
		// user database already populated, no onboarding allowed
		onboardFinished.Store(true)
	}

	return false, nil
}

func (a *Api) GetOnboardStatus(w http.ResponseWriter, r *http.Request) {
	can_ob, err := a.canOboard()
	Send(w, can_ob, err)
}

type OnboardRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Api) CreateOnboardUser(w http.ResponseWriter, r *http.Request) {
	can_ob, err := a.canOboard()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if !can_ob {
		SendError(w, 405, "onboarding not possible")
		return
	}

	obr := &OnboardRequest{}
	err = json.NewDecoder(r.Body).Decode(obr)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = auth.ValidateUsername(obr.Username)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("username invalid: %v", err))
		return
	}

	// create new password
	salt, verifier, err := auth.CreateAuth(obr.Username, obr.Password)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	auth_entry := &dbauth.AuthEntry{
		Name:     obr.Username,
		Password: auth.CreateDBPassword(salt, verifier),
	}
	// save to db
	err = a.app.DBContext.Auth.Create(auth_entry)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	for _, priv := range []string{"server", "interact", "privs"} {
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
		Username: obr.Username,
		Event:    "signup",
		Message:  fmt.Sprintf("User '%s' signed up successfully as admin", obr.Username),
	}, r)
}
