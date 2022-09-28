package web

import (
	"encoding/json"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

// fetch("api/bridge/lua", { method: "POST", body: JSON.stringify({ code: "return 123" })}).then(r => r.json()).then(res => console.log(res));
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
}
