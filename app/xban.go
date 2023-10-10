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
)

func (app *App) GetXBanDatabase() (*types.XBanDatabase, error) {
	data, err := os.ReadFile(path.Join(app.WorldDir, "xban.db"))
	if err != nil {
		return nil, fmt.Errorf("readfile error: %v", err)
	}

	xdb := &types.XBanDatabase{}
	err = json.Unmarshal(data, xdb)
	return xdb, err
}

// returns the xban entry from the on-disk xban db
func (app *App) GetOfflineXBanEntry(playername string) (*types.XBanEntry, error) {
	xdb, err := app.GetXBanDatabase()
	if err != nil {
		return nil, fmt.Errorf("xban db error: %v", err)
	}

	for _, e := range xdb.Entries {
		if e.Names[playername] {
			// entry found
			return e, nil
		}
	}

	// entry not found
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
