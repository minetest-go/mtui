package web

import (
	"fmt"
	"math/rand"
	"mtui/auth"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
	dbauth "github.com/minetest-go/mtdb/auth"
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// logs the specified user in as admin (jwt_key must be passed in the "key" query param)
// to be used in a provisioned environment for easy admin-onboarding
func (a *Api) AdminLogin(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	username := mux.Vars(r)["username"]

	if key != a.app.Config.JWTKey {
		SendError(w, 401, "invalid key")
		return
	}

	err := auth.ValidateUsername(username)
	if err != nil {
		SendError(w, 405, fmt.Sprintf("username invalid: %v", err))
		return
	}

	auth_entry, err := a.app.DBContext.Auth.GetByUsername(username)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("could not fetch auth entry: %v", err))
		return
	}

	if auth_entry == nil {
		// create new auth entry with random password
		salt, verifier, err := auth.CreateAuth(username, RandSeq(12))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		auth_entry = &dbauth.AuthEntry{
			Name:     username,
			Password: auth.CreateDBPassword(salt, verifier),
		}
		// save to db
		err = a.app.DBContext.Auth.Create(auth_entry)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	existing_privs, err := a.app.DBContext.Privs.GetByID(*auth_entry.ID)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("could not fetch priv entries: %v", err))
		return
	}

	for _, priv := range []string{"server", "interact", "privs", "ban"} {
		priv_exists := false
		for _, existing_priv := range existing_privs {
			if existing_priv.Privilege == priv {
				priv_exists = true
				break
			}
		}
		if priv_exists {
			// already there, skip
			continue
		}

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
		Username: username,
		Event:    "signup",
		Message:  fmt.Sprintf("User '%s' logged in successfully as admin", username),
	}, r)

	_, err = a.updateToken(w, *auth_entry.ID, username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// redirect to main page
	http.Redirect(w, r, fmt.Sprintf("%s#/help", a.app.Config.CookiePath), http.StatusSeeOther)
}
