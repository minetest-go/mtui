package web

import (
	"encoding/json"
	"fmt"
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
		SendError(w, 500, err)
		return
	}
	if !can_ob {
		SendError(w, 405, fmt.Errorf("onboarding not possible"))
		return
	}

	obr := &OnboardRequest{}
	err = json.NewDecoder(r.Body).Decode(obr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	_, err = a.app.CreateAdmin(obr.Username, obr.Password)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: obr.Username,
		Event:    "signup",
		Message:  fmt.Sprintf("User '%s' signed up successfully as admin", obr.Username),
	}, r)
}
