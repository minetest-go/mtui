package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *Api) GetLatestChatLogs(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	channel := mux.Vars(r)["channel"]
	list, err := a.app.Repos.ChatLogRepo.GetLatest(channel, 1000)
	Send(w, list, err)
}

func (a *Api) GetChatLogs(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	channel := vars["channel"]
	from, _ := strconv.ParseInt(vars["from"], 10, 64)
	to, _ := strconv.ParseInt(vars["to"], 10, 64)

	list, err := a.app.Repos.ChatLogRepo.Search(channel, from, to)
	Send(w, list, err)
}

func (a *Api) SendChat(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	msg := &command.ChatMessage{}
	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	msg.Name = c.Username

	err = a.app.Repos.ChatLogRepo.Insert(&types.ChatLog{
		Channel: msg.Channel,
		Name:    c.Username,
		Message: msg.Message,
	})
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = a.app.Bridge.SendCommand(command.COMMAND_CHAT_SEND, msg)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	log := &types.Log{
		Category: types.CategoryUI,
		Event:    "chat",
		Username: c.Username,
		Message:  fmt.Sprintf("'%s' writes '%s' in channel '%s'", c.Username, msg.Message, msg.Channel),
	}
	a.app.ResolveLogGeoIP(log, r)
	err = a.app.Repos.LogRepository.Insert(log)
	Send(w, msg, err)
}
