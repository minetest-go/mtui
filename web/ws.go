package web

import (
	"fmt"
	"mtui/eventbus"
	"mtui/types"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// event caching
var cachedEvents = map[eventbus.EventType]*eventbus.Event{}
var cache_lock = &sync.RWMutex{}

func sendEvent(wse *eventbus.Event, claims *types.Claims, conn *websocket.Conn) error {
	// check if a privilege is required for this event
	if wse.RequiredPriv != "" {
		if claims == nil {
			return nil
		}
		has_priv := false
		for _, priv := range claims.Privileges {
			if priv == wse.RequiredPriv {
				has_priv = true
				break
			}
		}
		if !has_priv {
			return nil
		}
	}
	return conn.WriteJSON(wse)
}

func (api *Api) WSCacheListener() {
	ch := make(chan *eventbus.Event, 1000)
	api.app.WSEvents.AddListener(ch)
	defer api.app.WSEvents.RemoveListener(ch)

	for wse := range ch {
		if wse.Cache {
			// cache event
			cache_lock.Lock()
			cachedEvents[wse.Type] = wse
			cache_lock.Unlock()
		}
	}
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

	// send cached events
	cache_lock.RLock()
	for _, ev := range cachedEvents {
		err = sendEvent(ev, claims, conn)
		if err != nil {
			fmt.Printf("WriteJSON: %s", err.Error())
			cache_lock.RUnlock()
			return
		}
	}
	cache_lock.RUnlock()

	// send live events
	for wse := range ch {
		err = sendEvent(wse, claims, conn)
		if err != nil {
			fmt.Printf("WriteJSON: %s", err.Error())
			return
		}
	}
}
