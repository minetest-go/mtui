package web

import (
	"mtui/types"
	"net/http"
)

func (a *Api) GetMails(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	list, err := a.app.Mail.List(claims.Username)
	Send(w, list, err)
}

func (a *Api) GetContacts(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c, err := a.app.Mail.Contacts(claims.Username)
	Send(w, c, err)
}
