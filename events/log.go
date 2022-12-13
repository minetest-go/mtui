package events

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/bridge"
	"mtui/db"
	"mtui/types"
)

func logLoop(lr *db.LogRepository, geoipresolver *app.GeoipResolver, ch chan *bridge.CommandResponse) {
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
				log.GeoCountry = &geoip.Country
				asn := int(geoip.ASN)
				log.GeoASN = &asn
			}
		}

		err = lr.Insert(log)
		if err != nil {
			fmt.Printf("DB error: %s\n", err.Error())
			return
		}
	}
}
