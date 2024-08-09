package web

import (
	"mtui/app"
	"mtui/auth"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) RequestJoinPassword(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	pw := app.RandSeq(16)
	salt, verifier, err := auth.CreateAuth(claims.Username, pw)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	dbpass := auth.CreateDBPassword(salt, verifier)

	req := &command.SetJoinPasswordRequest{
		Username: claims.Username,
		Password: dbpass,
	}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_SET_JOIN_PASSWORD, req, nil, time.Second*2)
	Send(w, pw, err)
}
