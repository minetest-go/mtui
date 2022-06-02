package web

import (
	"fmt"
	"mtui/eventbus"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (api *Api) Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	ch := make(chan *eventbus.Event, 1000)
	api.app.WSEvents.AddListener(ch)
	defer api.app.WSEvents.RemoveListener(ch)

	for wse := range ch {
		err := conn.WriteJSON(wse)
		if err != nil {
			fmt.Printf("WriteJSON: %s", err.Error())
			return
		}
	}
}
