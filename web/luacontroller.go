package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) GetLuacontroller(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerGetProgramRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if req.Pos == nil {
		SendError(w, 500, "position is nil")
		return
	}
	req.Playername = claims.Username

	a.app.CreateUILogEntry(&types.Log{
		Category: types.CategoryUI,
		Event:    "mesecons",
		PosX:     &req.Pos.X,
		PosY:     &req.Pos.Y,
		PosZ:     &req.Pos.Z,
		Message:  fmt.Sprintf("user '%s' requests code from luacontroller", claims.Username),
	}, r)

	resp := &command.LuaControllerGetProgramResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_GETPROGRAM_LUACONTROLLER, req, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) SetLuacontroller(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerSetProgramRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if req.Pos == nil {
		SendError(w, 500, "position is nil")
		return
	}
	req.Playername = claims.Username

	a.app.CreateUILogEntry(&types.Log{
		Category:   types.CategoryUI,
		Event:      "mesecons",
		PosX:       &req.Pos.X,
		PosY:       &req.Pos.Y,
		PosZ:       &req.Pos.Z,
		Message:    fmt.Sprintf("user '%s' writes code to luacontroller", claims.Username),
		Attachment: []byte(req.Code),
	}, r)

	resp := &command.LuaControllerSetProgramResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_SETPROGRAM_LUACONTROLLER, req, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) LuacontrollerDigilineSend(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerDigilineSendRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if req.Pos == nil {
		SendError(w, 500, "position is nil")
		return
	}
	if len(req.Channel) > 64 || len(req.Message) > 128 {
		SendError(w, 500, "channel/message length exceeded")
		return
	}
	req.Playername = claims.Username

	a.app.CreateUILogEntry(&types.Log{
		Category: types.CategoryUI,
		Event:    "mesecons",
		PosX:     &req.Pos.X,
		PosY:     &req.Pos.Y,
		PosZ:     &req.Pos.Z,
		Message:  fmt.Sprintf("user '%s' sends digiline message '%s' on channel '%s'", claims.Username, req.Message, req.Channel),
	}, r)

	resp := &command.LuaControllerDigilineSendResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_LUACONTROLLER_DIGLINE_SEND, req, resp, time.Second*2)
	Send(w, resp, err)
}
