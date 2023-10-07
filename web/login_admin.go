package web

import (
	"fmt"
	"math/rand"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
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

	// new admin with random password
	auth_entry, err := a.app.CreateAdmin(username, RandSeq(12))
	if err != nil {
		SendError(w, 500, err.Error())
		return
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
