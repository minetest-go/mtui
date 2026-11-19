package web

import (
	"encoding/json"
	"errors"
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

// findMail loads the user's mail entry and the message with the id from the request.
// On failure the error response is written and ok is false
func (a *Api) findMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) (e *mail.PlayerEntry, id string, m *mail.Message, ok bool) {
	id = mux.Vars(r)["id"]
	e, err := a.app.Mail.GetEntry(claims.Username)
	if err != nil {
		SendError(w, 500, err)
		return nil, id, nil, false
	}
	m = e.FindMessage(id)
	if m == nil {
		SendError(w, 404, fmt.Errorf("mail not found"))
		return nil, id, nil, false
	}
	return e, id, m, true
}

func (a *Api) markMail(w http.ResponseWriter, r *http.Request, claims *types.Claims, read bool) {
	e, _, m, ok := a.findMail(w, r, claims)
	if !ok {
		return
	}
	m.Read = read
	err := a.app.Mail.SetEntry(claims.Username, e)
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

// errInvalidRecipient lists the recipients that don't exist
type errInvalidRecipient []string

func (e errInvalidRecipient) Error() string {
	return fmt.Sprintf("invalid recipient: '%s'", strings.Join(e, "', '"))
}

// parsePlayerList splits a recipient field on commas and spaces and removes
// case-insensitive duplicates, like mail.parse_player_list in the mail mod
func parsePlayerList(field string) []string {
	seen := map[string]bool{}
	list := []string{}
	for _, name := range strings.FieldsFunc(field, func(r rune) bool { return r == ',' || r == ' ' }) {
		lower := strings.ToLower(name)
		if !seen[lower] {
			seen[lower] = true
			list = append(list, name)
		}
	}
	return list
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

// sendMail validates the recipients and delivers the message from the user to every
// recipient in to, cc and bcc (same storage layout as mail.send in the mail mod)
func (a *Api) sendMail(claims *types.Claims, m *mail.Message) error {
	// populate fixed fields
	m.ID = uuid.NewString()
	m.From = claims.Username
	m.Time = float64(time.Now().Unix())

	// normalize the recipient fields and collect the unique recipients
	recipients := []string{}
	seen := map[string]bool{}
	normalize := func(field string) string {
		list := parsePlayerList(field)
		for _, name := range list {
			if !seen[name] {
				seen[name] = true
				recipients = append(recipients, name)
			}
		}
		return strings.Join(list, ", ")
	}
	normalizeOptional := func(field *string) *string {
		if field == nil {
			return nil
		}
		v := normalize(*field)
		if v == "" {
			return nil
		}
		return &v
	}
	m.To = normalize(m.To)
	m.CC = normalizeOptional(m.CC)
	m.BCC = normalizeOptional(m.BCC)

	if len(recipients) == 0 {
		return errors.New("no recipient")
	}

	undeliverable := errInvalidRecipient{}
	for _, rcpt := range recipients {
		valid, err := a.isValidRecipient(rcpt)
		if err != nil {
			return err
		}
		if !valid {
			undeliverable = append(undeliverable, rcpt)
		}
	}
	if len(undeliverable) > 0 {
		return undeliverable
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

	// insert into every receivers inbox
	for _, rcpt := range recipients {
		e, err = a.app.Mail.GetEntry(rcpt)
		if err != nil {
			return err
		}
		e.Inbox = append(e.Inbox, m)
		err = a.app.Mail.SetEntry(rcpt, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Api) DeleteMail(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	e, id, m, ok := a.findMail(w, r, claims)
	if !ok {
		return
	}
	e.RemoveMessage(id)
	err := a.app.Mail.SetEntry(claims.Username, e)
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
