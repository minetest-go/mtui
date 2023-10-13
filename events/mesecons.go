package events

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"time"

	"github.com/sirupsen/logrus"
)

func meseconsRegister(a *app.App, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		m := &command.MeseconsRegister{}
		err := json.Unmarshal(cmd.Data, m)
		if err != nil {
			logrus.WithError(err).Error("mesecons register payload error")
			continue
		}

		me, err := a.Repos.MeseconsRepo.GetByPoskey(m.Pos.String())
		if err != nil {
			logrus.WithError(err).Error("get by poskey error")
			continue
		}

		if me == nil {
			// new entry
			me = &types.Mesecons{
				PosKey: m.Pos.String(),
				X:      int(m.Pos.X),
				Y:      int(m.Pos.Y),
				Z:      int(m.Pos.Z),
				Name:   fmt.Sprintf("registered '%s'", m.Nodename),
			}
		}

		// update with incoming fields
		me.PlayerName = m.Playername
		me.NodeName = m.Nodename
		me.LastModified = time.Now().UnixMilli()

		// get next order id
		list, err := a.Repos.MeseconsRepo.GetByPlayerName(m.Playername)
		if err != nil {
			logrus.WithError(err).Error("mesecons get error")
			continue
		}
		if len(list) > 0 {
			me.OrderID = list[len(list)-1].OrderID + 1
		}

		err = a.Repos.MeseconsRepo.Save(me)
		if err != nil {
			logrus.WithError(err).Error("mesecons save error")
			continue
		}
	}
}

func meseconsEvent(a *app.App, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		ev := &command.MeseconsEvent{}
		err := json.Unmarshal(cmd.Data, ev)
		if err != nil {
			logrus.WithError(err).Error("mesecons register payload error")
			continue
		}

		me, err := a.Repos.MeseconsRepo.GetByPoskey(ev.Pos.String())
		if err != nil {
			logrus.WithError(err).Error("get by poskey error")
			continue
		}
		if me == nil {
			// no such entry
			continue
		}

		// update state
		me.State = string(ev.State)
		me.NodeName = ev.Nodename
		me.LastModified = time.Now().UnixMilli()
		err = a.Repos.MeseconsRepo.Save(me)
		if err != nil {
			logrus.WithError(err).Error("mesecons save error")
			continue
		}
	}
}
