package web

import (
	"encoding/json"
	"mtui/mail"
	"mtui/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) GetMails(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	list, err := a.app.Mail.GetMessages(claims.Username)
	Send(w, list, err)
}

func (a *Api) CheckValidRecipient(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	auth_entry, err := a.app.DBContext.Auth.GetByUsername(recipient)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth_entry == nil {
		SendError(w, 404, "Recipient player does not exist")
		return
	}
}

func (a *Api) DeleteMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	sender := vars["sender"]
	time, err := strconv.ParseInt(vars["time"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := a.app.Mail.GetMessages(claims.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// filter removed mail
	new_list := make([]*mail.Message, 0)
	for _, msg := range list {
		if msg.Sender != sender || int64(msg.Time) != time {
			new_list = append(new_list, msg)
		}
	}

	err = a.app.Mail.SetMessages(claims.Username, new_list)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (a *Api) MarkRead(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	sender := vars["sender"]
	time, err := strconv.ParseInt(vars["time"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := a.app.Mail.GetMessages(claims.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	for _, msg := range list {
		if msg.Sender == sender && int64(msg.Time) == time {
			msg.Unread = false
		}
	}

	err = a.app.Mail.SetMessages(claims.Username, list)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (a *Api) SendMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	auth_entry, err := a.app.DBContext.Auth.GetByUsername(recipient)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth_entry == nil {
		SendError(w, 404, "Recipient player does not exist")
		return
	}

	msg := &mail.Message{}
	err = json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// set known fields
	msg.Sender = claims.Username
	msg.Time = float64(time.Now().Unix())
	msg.Unread = true

	recipient_mails, err := a.app.Mail.GetMessages(recipient)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	recipient_mails = append(recipient_mails, msg)
	err = a.app.Mail.SetMessages(recipient, recipient_mails)

	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (a *Api) GetContacts(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	c, err := a.app.Mail.GetContacts(claims.Username)
	Send(w, c, err)
}
