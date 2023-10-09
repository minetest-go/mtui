package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"net/http"

	"github.com/dchest/captcha"
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

	_, err = a.app.CreateUser(sr.Username, sr.Password, false, []string{"shout", "interact"})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: sr.Username,
		Event:    "signup",
		Message:  fmt.Sprintf("User '%s' signed up successfully", sr.Username),
	}, r)
}
