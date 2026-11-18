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

	resp, err := a.executeChatcommand(claims, req.Command, r)
	Send(w, resp, err)
}

// executeChatcommand runs the chatcommand ingame as the user
func (a *Api) executeChatcommand(claims *types.Claims, cmd string, r *http.Request) (*command.ExecuteChatCommandResponse, error) {
	req := &command.ExecuteChatCommandRequest{
		Playername: claims.Username,
		Command:    cmd,
	}

	resp := &command.ExecuteChatCommandResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_CHATCMD, req, resp, time.Second*5)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "chatcommand",
		Message:  fmt.Sprintf("User '%s' executes the chatcommand: '%s'", claims.Username, cmd),
	}, r)

	return resp, err
}
