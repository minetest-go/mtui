package web

import (
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"net/http"
)

var features = types.FeaturesCommand{}

func (a *Api) FeatureListener(c chan *bridge.Command) {
	for {
		cmd := <-c
		o, err := types.ParseCommand(cmd)
		if err != nil {
			fmt.Printf("Feature-listener-error: %s\n", err.Error())
			continue
		}

		switch payload := o.(type) {
		case *types.FeaturesCommand:
			features = *payload
		}
	}
}

func (a *Api) GetFeatures(w http.ResponseWriter, r *http.Request) {
	SendJson(w, features)
}
