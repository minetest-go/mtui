package web

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/types"
	"net/http"
)

func (a *Api) GetFeatures(w http.ResponseWriter, r *http.Request) {
	feature_map := make(map[string]*app.AvailableFeature)

	available_features, err := app.GetAvailableFeatures()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	for _, avavailable_feature := range available_features {
		feature_map[avavailable_feature.Name] = avavailable_feature
	}

	list, err := a.app.Repos.FeatureRepository.GetAll()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	for _, feature := range list {
		f := feature_map[feature.Name]
		if f != nil {
			f.Enabled = feature.Enabled
		}
	}

	SendJson(w, feature_map)
}

func (a *Api) SetFeature(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	feature := &types.Feature{}
	err := json.NewDecoder(r.Body).Decode(feature)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = a.app.Repos.FeatureRepository.Set(feature)
	Send(w, feature, err)

	// create log entry
	action := "disables"
	if feature.Enabled {
		action = "enables"
	}
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "feature",
		Message:  fmt.Sprintf("User '%s' %s the feature '%s'", claims.Username, action, feature.Name),
	}, r)
}
