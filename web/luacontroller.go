package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

// authorizeMesecons checks that the position is set and owned by the user.
// On failure the error response is written and false is returned
func (a *Api) authorizeMesecons(w http.ResponseWriter, pos *types.Pos, claims *types.Claims) bool {
	if pos == nil {
		SendError(w, 500, fmt.Errorf("position is nil"))
		return false
	}

	m, err := a.app.Repos.MeseconsRepo.GetByPoskey(pos.String())
	if err != nil {
		SendError(w, 500, err)
		return false
	}
	if m == nil || m.PlayerName != claims.Username {
		SendError(w, 403, fmt.Errorf("unauthorized"))
		return false
	}
	return true
}

func (a *Api) logMesecons(r *http.Request, claims *types.Claims, pos *types.Pos, attachment []byte, format string, args ...any) {
	a.app.CreateUILogEntry(&types.Log{
		Category:   types.CategoryUI,
		Event:      "mesecons",
		PosX:       &pos.X,
		PosY:       &pos.Y,
		PosZ:       &pos.Z,
		Username:   claims.Username,
		Message:    fmt.Sprintf(format, args...),
		Attachment: attachment,
	}, r)
}

func (a *Api) GetLuacontroller(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerGetProgramRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, 500, err)
		return
	}
	if !a.authorizeMesecons(w, req.Pos, claims) {
		return
	}

	req.Playername = claims.Username
	a.logMesecons(r, claims, req.Pos, nil, "user '%s' requests code from luacontroller", claims.Username)

	resp := &command.LuaControllerGetProgramResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_GETPROGRAM_LUACONTROLLER, req, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) SetLuacontroller(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerSetProgramRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, 500, err)
		return
	}
	if !a.authorizeMesecons(w, req.Pos, claims) {
		return
	}

	req.Playername = claims.Username
	a.logMesecons(r, claims, req.Pos, []byte(req.Code), "user '%s' writes code to luacontroller", claims.Username)

	resp := &command.LuaControllerSetProgramResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_SETPROGRAM_LUACONTROLLER, req, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) LuacontrollerDigilineSend(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.LuaControllerDigilineSendRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, 500, err)
		return
	}
	if req.Pos != nil && (len(req.Channel) > 64 || len(req.Message) > 128) {
		SendError(w, 500, fmt.Errorf("channel/message length exceeded"))
		return
	}
	if !a.authorizeMesecons(w, req.Pos, claims) {
		return
	}

	req.Playername = claims.Username
	a.logMesecons(r, claims, req.Pos, nil, "user '%s' sends digiline message '%s' on channel '%s'", claims.Username, req.Message, req.Channel)

	resp := &command.LuaControllerDigilineSendResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_LUACONTROLLER_DIGLINE_SEND, req, resp, time.Second*2)
	Send(w, resp, err)
}
