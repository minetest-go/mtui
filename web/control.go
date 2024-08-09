package web

import (
	"encoding/json"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"
)

func (a *Api) GetControlsMetadata(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	resp := &command.GetControlsMetadataResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_GET_CONTROLS_METADATA, nil, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) GetControlsValues(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	resp := &command.GetControlsValuesResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_GET_CONTROLS_VALUES, nil, resp, time.Second*2)
	Send(w, resp, err)
}

func (a *Api) SetControl(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.SetControlRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	success := false
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_SET_CONTROL, req, &success, time.Second*2)
	Send(w, success, err)
}
