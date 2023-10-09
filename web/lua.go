package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) ExecuteLua(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	resp := &command.LuaResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	Send(w, resp, err)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "lua",
		Message:  fmt.Sprintf("User '%s' executes the lua-code: '%s'", claims.Username, req.Code),
	}, r)
}
