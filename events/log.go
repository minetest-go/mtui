package events

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/bridge"
	"mtui/types"
)

func logLoop(a *app.App, geoipresolver *app.GeoipResolver, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		log := &types.Log{}
		err := json.Unmarshal(cmd.Data, log)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}

		if log.IPAddress != nil {
			geoip := geoipresolver.Resolve(*log.IPAddress)
			if geoip != nil {
				log.GeoCity = &geoip.City
				log.GeoCountry = &geoip.ISOCountry
				log.GeoASN = &geoip.ASN
			}
		}

		err = a.Repos.LogRepository.Insert(log)
		if err != nil {
			fmt.Printf("DB error: %s\n", err.Error())
			return
		}
	}
}
