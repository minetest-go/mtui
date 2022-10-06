package web

import (
	"encoding/json"
	"mtui/types"
	"net/http"
)

func (a *Api) GetFeatures(w http.ResponseWriter, r *http.Request) {
	feature_map := make(map[string]bool)

	list, err := a.app.Repos.FeatureRepository.GetAll()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	for _, feature := range list {
		feature_map[feature.Name] = feature.Enabled
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
}
