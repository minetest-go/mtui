package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) GetATMBalance(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	name := vars["name"]

	entry, err := a.app.DBContext.ModStorage.Get("atm", []byte(fmt.Sprintf("balance_%s", name)))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	balance := 0

	if entry != nil {
		v, err := strconv.ParseInt(string(entry.Value), 10, 64)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		balance = int(v)
	}

	Send(w, map[string]int{"balance": balance}, nil)
}

func (a *Api) ATMTransfer(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	req := &command.ATMTransferRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// set logged in user as source
	req.Source = claims.Username

	resp := &command.ATMTransferResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_ATM_TRANSFER, req, resp, time.Second*2)
	Send(w, resp, err)

	a.app.CreateUILogEntry(&types.Log{
		Category: types.CategoryUI,
		Event:    "atm",
		Username: claims.Username,
		Message:  fmt.Sprintf("User '%s' transfers $ %d to user '%s'", claims.Username, req.Amount, req.Target),
	}, r)
}
