package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SendLuaResponse(w http.ResponseWriter, err error, lr *command.LuaResponse) {
	if err != nil {
		SendError(w, 500, err)
	} else if !lr.Success {
		SendError(w, 500, errors.New(lr.Message))
	} else {
		Send(w, lr.Result, nil)
	}
}

func (a *Api) GetBanDBStatus(w http.ResponseWriter, r *http.Request) {
	req := &command.LuaRequest{
		Code: `
			local banned = 0
			local total = 0
			for _, entry in ipairs(xban.db) do
				total = total + 1
				if entry.banned then
					banned = banned + 1
				end
			end
			return { total = total, banned = banned }
		`,
	}
	resp := &command.LuaResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	SendLuaResponse(w, err, resp)
}

func (a *Api) GetBannedRecords(w http.ResponseWriter, r *http.Request) {
	result, err := a.getBannedRecords()
	Send(w, result, err)
}

func (a *Api) getBannedRecords() (any, error) {
	req := &command.LuaRequest{
		Code: `
			local banned = {}
			for _, entry in ipairs(xban.db) do
				if entry.banned then
					table.insert(banned, entry)
				end
			end
			return banned
		`,
	}
	return a.executeLuaResult(req.Code)
}

func (a *Api) GetBanRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	e, err := a.app.GetOnlineXBanEntry(vars["playername"])
	Send(w, e, err)
}

func (a *Api) BanPlayer(w http.ResponseWriter, r *http.Request) {
	a.xbanAction(w, r, func(c *types.Claims, banr *types.XBanRequest) (any, error) {
		return a.banPlayer(c, banr.Playername, banr.Reason, nil)
	})
}

func (a *Api) TempBanPlayer(w http.ResponseWriter, r *http.Request) {
	a.xbanAction(w, r, func(c *types.Claims, banr *types.XBanRequest) (any, error) {
		seconds := int64(banr.Time)
		return a.banPlayer(c, banr.Playername, banr.Reason, &seconds)
	})
}

func (a *Api) UnbanPlayer(w http.ResponseWriter, r *http.Request) {
	a.xbanAction(w, r, func(c *types.Claims, banr *types.XBanRequest) (any, error) {
		return a.unbanPlayer(c, banr.Playername)
	})
}

func (a *Api) xbanAction(w http.ResponseWriter, r *http.Request, fn func(*types.Claims, *types.XBanRequest) (any, error)) {
	claims, err := a.GetClaims(r)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	banr := &types.XBanRequest{}
	err = json.NewDecoder(r.Body).Decode(banr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	res, err := fn(claims, banr)
	Send(w, res, err)
}

// banPlayer bans the player with the user as source, a nil duration is a permanent ban
func (a *Api) banPlayer(claims *types.Claims, playername, reason string, seconds *int64) (any, error) {
	duration := "nil"
	if seconds != nil {
		duration = fmt.Sprintf("%d", *seconds)
	}
	return a.executeLuaResult(fmt.Sprintf("return xban.ban_player('%s', '%s', %s, '%s')",
		bridge.SanitizeLuaString(playername),
		bridge.SanitizeLuaString(claims.Username),
		duration,
		bridge.SanitizeLuaString(reason),
	))
}

func (a *Api) unbanPlayer(claims *types.Claims, playername string) (any, error) {
	return a.executeLuaResult(fmt.Sprintf("return xban.unban_player('%s', '%s')",
		bridge.SanitizeLuaString(playername),
		bridge.SanitizeLuaString(claims.Username),
	))
}

// executeLuaResult runs the code in the engine and returns the result or the lua error
func (a *Api) executeLuaResult(code string) (any, error) {
	req := &command.LuaRequest{Code: code}
	resp := &command.LuaResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, errors.New(resp.Message)
	}
	return resp.Result, nil
}

func (a *Api) CleanupBanDB(w http.ResponseWriter, r *http.Request) {
	req := &command.LuaRequest{
		Code: `
			local db = xban.db
			local old_count = #db
			local i = 1
			while i <= #db do
				if not db[i].banned then
					-- not banned, remove from db
					table.remove(db, i)
				else
					-- banned, hold entry back
					i = i + 1
				end
			end
			return {
				removed = (old_count - #db),
				retained = #db
			}
		`,
	}
	resp := &command.LuaResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	SendLuaResponse(w, err, resp)
}
