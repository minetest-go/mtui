package web

import (
	"encoding/json"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) GetMeseconsControls(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	list, err := a.app.Repos.MeseconsRepo.GetByPlayerName(claims.Username)
	Send(w, list, err)
}

func (a *Api) SetMeseconsControl(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	user_mesecon := &types.Mesecons{}
	err := json.NewDecoder(r.Body).Decode(user_mesecon)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	mesecon, err := a.app.Repos.MeseconsRepo.GetByPoskey(user_mesecon.PosKey)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if mesecon == nil {
		SendError(w, 404, "not found")
		return
	}
	if mesecon.PlayerName != claims.Username {
		SendError(w, 403, "unauthorized")
		return
	}

	// update user provided fields
	mesecon.Name = user_mesecon.Name
	mesecon.State = user_mesecon.State
	mesecon.LastModified = time.Now().UnixMilli()
	err = a.app.Repos.MeseconsRepo.Save(mesecon)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	cmd_req := &command.MeseconsSetRequest{}
	cmd_resp := &command.MeseconsSetResponse{}

	err = a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_SET, cmd_req, cmd_resp, time.Second)
	Send(w, cmd_resp, err)
}