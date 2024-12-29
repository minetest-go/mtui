package app

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"os"
	"path"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
)

func (app *App) GetXBanDatabase() (*types.XBanDatabase, error) {
	data, err := os.ReadFile(path.Join(app.WorldDir, "xban.db"))
	if err != nil {
		// file not found or other file-related error
		return nil, nil
	}

	xdb := &types.XBanDatabase{}
	err = json.Unmarshal(data, xdb)
	return xdb, err
}

// returns the xban entry from the on-disk xban db
func (app *App) GetOfflineXBanEntry(playername string) (*types.XBanEntry, error) {
	cached_entry, ok := app.offline_xban_cache.Get(playername)
	if ok {
		return cached_entry, nil
	}

	xdb, err := app.GetXBanDatabase()
	if err != nil {
		return nil, fmt.Errorf("xban db error: %v", err)
	}
	if xdb == nil {
		return nil, nil
	}

	for _, e := range xdb.Entries {
		if e.Names[playername] {
			// entry found
			app.offline_xban_cache.Set(playername, e, cache.WithExpiration(time.Minute))
			return e, nil
		}
	}

	// entry not found
	app.offline_xban_cache.Set(playername, nil, cache.WithExpiration(time.Minute))
	return nil, nil
}

// returns the xban entry directly from the live server
func (app *App) GetOnlineXBanEntry(playername string) (*types.XBanEntry, error) {
	e := &types.XBanEntry{}

	req := &command.LuaRequest{
		Code: fmt.Sprintf("return xban.find_entry('%s')", bridge.SanitizeLuaString(playername)),
	}
	resp := &command.LuaResponse{}

	err := app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*1)
	if err != nil {
		return nil, fmt.Errorf("could not get xban entry: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("non-success resut: '%s'", resp.Message)
	}

	if resp.Result == nil {
		return nil, nil
	}

	err = json.Unmarshal(resp.Result, e)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}

	return e, nil
}
