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
	claims, err := GetClaims(r)
	if err != nil && err != err_unauthorized {
		return
	}

	defer conn.Close()
	ch := make(chan *eventbus.Event, 1000)
	api.app.WSEvents.AddListener(ch)
	defer api.app.WSEvents.RemoveListener(ch)

	for wse := range ch {
		// check if a privilege is required for this event
		if wse.RequiredPriv != "" {
			if claims == nil {
				continue
			}
			has_priv := false
			for _, priv := range claims.Privileges {
				if priv == wse.RequiredPriv {
					has_priv = true
					break
				}
			}
			if !has_priv {
				continue
			}
		}
		err := conn.WriteJSON(wse)
		if err != nil {
			fmt.Printf("WriteJSON: %s", err.Error())
			return
		}
	}
}
