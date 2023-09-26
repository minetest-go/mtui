package web

import (
	"encoding/json"
	"fmt"
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

	is_superuser := claims.HasPriv("server") || claims.HasPriv("password")

	// check username
	if req.Username != claims.Username && !is_superuser {
		// username does not match and "server" or "password" not found (not a superuser)
		SendError(w, 500, "username mismatch and not a superuser")
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

	if !is_superuser {
		// check old password (not a super-user)

		// legacy password first
		legacy_ok := auth.VerifyLegacyPassword(req.Username, req.OldPassword, auth_entry.Password)
		if !legacy_ok {
			// SRP fallback
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
		}
	}

	// create new password
	salt, verifier, err := auth.CreateAuth(req.Username, req.NewPassword)
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

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "password",
		Message:  fmt.Sprintf("User '%s' set the passowrd of %s", claims.Username, req.Username),
	}, r)
}
