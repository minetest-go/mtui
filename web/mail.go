package web

import (
	"encoding/json"
	"fmt"
	"mtui/mail"
	"mtui/types"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetContacts(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	e, err := a.app.Mail.GetEntry(claims.Username)
	Send(w, e.Contacts, err)
}

func (a *Api) GetInbox(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	e, err := a.app.Mail.GetEntry(claims.Username)
	Send(w, e.Inbox, err)
}

func (a *Api) GetOutbox(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	e, err := a.app.Mail.GetEntry(claims.Username)
	Send(w, e.Outbox, err)
}

func (a *Api) GetDrafts(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	e, err := a.app.Mail.GetEntry(claims.Username)
	Send(w, e.Drafts, err)
}

func (a *Api) markMail(w http.ResponseWriter, r *http.Request, claims *types.Claims, read bool) {
	vars := mux.Vars(r)
	id := vars["id"]
	e, err := a.app.Mail.GetEntry(claims.Username)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	m := e.FindMessage(id)
	if m == nil {
		SendError(w, 404, fmt.Errorf("mail not found"))
		return
	}
	m.Read = read
	err = a.app.Mail.SetEntry(claims.Username, e)
	Send(w, m, err)
}

func (a *Api) MarkRead(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	a.markMail(w, r, claims, true)
}

func (a *Api) MarkUnread(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	a.markMail(w, r, claims, false)
}

func (a *Api) isValidRecipient(recipient string) (bool, error) {
	auth_entry, err := a.app.DBContext.Auth.GetByUsername(recipient)
	if err != nil {
		return false, err
	}
	return auth_entry != nil, nil
}

type errInvalidRecipient string

func (e errInvalidRecipient) Error() string {
	return fmt.Sprintf("invalid recipient: '%s'", string(e))
}

func (a *Api) checkRecipient(recipient string) error {
	v, err := a.isValidRecipient(recipient)
	if err != nil {
		return err
	}
	if !v {
		return errInvalidRecipient(recipient)
	}
	return nil
}

func (a *Api) SendMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &mail.Message{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	err = a.sendMail(claims, m)
	if _, ok := err.(errInvalidRecipient); ok {
		SendError(w, 404, err)
		return
	}
	Send(w, m, err)
}

// sendMail validates the recipients and delivers the message from the user to the recipient
func (a *Api) sendMail(claims *types.Claims, m *mail.Message) error {
	// populate fixed fields
	m.ID = uuid.NewString()
	m.From = claims.Username
	m.Time = float64(time.Now().Unix())

	recipients := []string{m.To}
	if m.CC != nil {
		recipients = append(recipients, strings.Split(*m.CC, ",")...)
	}
	if m.BCC != nil {
		recipients = append(recipients, strings.Split(*m.BCC, ",")...)
	}
	for _, rcpt := range recipients {
		if err := a.checkRecipient(rcpt); err != nil {
			return err
		}
	}

	// insert into senders outbox
	e, err := a.app.Mail.GetEntry(claims.Username)
	if err != nil {
		return err
	}
	e.Outbox = append(e.Outbox, m)
	err = a.app.Mail.SetEntry(claims.Username, e)
	if err != nil {
		return err
	}

	// insert into receivers inbox
	e, err = a.app.Mail.GetEntry(m.To)
	if err != nil {
		return err
	}
	e.Inbox = append(e.Inbox, m)
	return a.app.Mail.SetEntry(m.To, e)
}

func (a *Api) DeleteMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	e, err := a.app.Mail.GetEntry(claims.Username)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	m := e.FindMessage(id)
	if m == nil {
		SendError(w, 404, fmt.Errorf("mail not found"))
		return
	}
	e.RemoveMessage(id)
	err = a.app.Mail.SetEntry(claims.Username, e)
	Send(w, m, err)

}

func (a *Api) CheckValidRecipient(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	v, err := a.isValidRecipient(recipient)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if !v {
		SendError(w, 404, fmt.Errorf("recipient player does not exist"))
		return
	}
}
