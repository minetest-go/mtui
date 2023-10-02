package web

import (
	"encoding/json"
	"mtui/types"
	"mtui/types/command"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetLatestChatLogs(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	channel := mux.Vars(r)["channel"]
	list, err := a.app.Repos.ChatLogRepo.GetLatest(channel, 1000)
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
	Send(w, msg, err)
}
