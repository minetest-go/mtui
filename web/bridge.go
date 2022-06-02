package web

import (
	"encoding/json"
	"mtadmin/types"
	"net/http"
	"time"
)

func (api *Api) BridgeRx(w http.ResponseWriter, r *http.Request) {
	cmd := &types.Command{}
	err := json.NewDecoder(r.Body).Decode(cmd)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	select {
	case api.rx_cmds <- cmd:
	default:
	}
}

// collects commands for a certain amount of time
func collectCommands(ch chan *types.Command, delay time.Duration) []*types.Command {
	cmd_list := make([]*types.Command, 0)
	then := time.Now().Add(delay)
	t := time.NewTimer(delay)

	for {
		select {
		case cmd := <-ch:
			// command received, add to list
			cmd_list = append(cmd_list, cmd)
		case <-t.C:
			// nothing received in the entire duration time
		}

		if time.Now().After(then) {
			// time is up, return collected commands, if any
			return cmd_list
		}
	}
}

func (api *Api) BridgeTx(w http.ResponseWriter, r *http.Request) {
	then := time.Now().Add(20 * time.Second)
	cmds := make([]*types.Command, 0)
	for {
		// collect commands for at least 100ms
		cmds = collectCommands(api.tx_cmds, 100*time.Millisecond)

		if len(cmds) > 0 {
			// commands received, return them
			break
		}

		if time.Now().After(then) {
			// time is up and not commands received
			break
		}
	}

	Send(w, cmds, nil)
}
