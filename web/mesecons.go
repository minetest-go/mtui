package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	mesecon.OrderID = user_mesecon.OrderID
	mesecon.Category = user_mesecon.Category
	mesecon.LastModified = time.Now().UnixMilli()
	err = a.app.Repos.MeseconsRepo.Save(mesecon)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	cmd_req := &command.MeseconsSetRequest{
		Pos: &types.Pos{
			X: types.JsonInt(mesecon.X),
			Y: types.JsonInt(mesecon.Y),
			Z: types.JsonInt(mesecon.Z),
		},
		State:    command.State(mesecon.State),
		Nodename: mesecon.NodeName,
	}
	cmd_resp := &command.MeseconsSetResponse{}

	a.app.CreateUILogEntry(&types.Log{
		Category: types.CategoryUI,
		Event:    "mesecons",
		PosX:     &cmd_req.Pos.X,
		PosY:     &cmd_req.Pos.Y,
		PosZ:     &cmd_req.Pos.Z,
		Message:  fmt.Sprintf("user '%s' sets mesecons to '%s'", claims.Username, user_mesecon.State),
	}, r)

	err = a.app.Bridge.ExecuteCommand(command.COMMAND_MESECONS_SET, cmd_req, cmd_resp, time.Second)
	Send(w, cmd_resp, err)
}

func (a *Api) DeleteMeseconsControl(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	poskey := vars["poskey"]

	m, err := a.app.Repos.MeseconsRepo.GetByPoskey(poskey)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if m.PlayerName != claims.Username {
		SendError(w, 403, "unauthorized")
		return
	}

	err = a.app.Repos.MeseconsRepo.Remove(poskey)
	Send(w, true, err)
}
