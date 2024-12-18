package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) ExecuteChatcommand(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.ExecuteChatCommandRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// assign playername from claims
	req.Playername = claims.Username

	resp := &command.ExecuteChatCommandResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_CHATCMD, req, resp, time.Second*5)
	Send(w, resp, err)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "chatcommand",
		Message:  fmt.Sprintf("User '%s' executes the chatcommand: '%s'", claims.Username, req.Command),
	}, r)

}
