package web

import "net/http"

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {

}
