package web

import (
	"encoding/json"
	"mtui/auth"
	"mtui/types"
	"net/http"
)

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (a *Api) ChangePassword(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &ChangePasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// check username
	if req.Username != claims.Username {
		// username does not match
		SendError(w, 500, "username mismatch")
		return
	}

	// fetch entry from db
	auth_entry, err := a.app.DBContext.Auth.GetByUsername(req.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth_entry == nil {
		SendError(w, 404, "not found")
		return
	}

	// check old password
	salt, verifier, err := auth.ParseDBPassword(auth_entry.Password)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	ok, err := auth.VerifyAuth(req.Username, req.OldPassword, salt, verifier)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if !ok {
		SendError(w, 401, "unauthorized")
		return
	}

	// create new password
	salt, verifier, err = auth.CreateAuth(req.Username, req.NewPassword)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// save to db
	auth_entry.Password = auth.CreateDBPassword(salt, verifier)
	err = a.app.DBContext.Auth.Update(auth_entry)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
