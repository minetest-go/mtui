package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"mtui/db"
	"mtui/types"
)

type AvailableFeature struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Mods         []string `json:"mods"`
	Experimental bool     `json:"experimental"`
	Enabled      bool     `json:"enabled"`
}

//go:embed features.json
var fs embed.FS

func GetAvailableFeatures() ([]*AvailableFeature, error) {
	feature_json, err := fs.ReadFile("features.json")
	if err != nil {
		return nil, err
	}
	available_features := make([]*AvailableFeature, 0)
	err = json.Unmarshal(feature_json, &available_features)
	return available_features, err
}

func PopulateFeatures(repo *db.FeatureRepository, enabled_features []string) error {

	available_features, err := GetAvailableFeatures()
	if err != nil {
		return err
	}

	for _, available_feature := range available_features {
		feature, err := repo.GetByName(available_feature.Name)
		if err != nil {
			return err
		}
		if feature == nil {
			// check if the feature was enabled in the env-vars
			for _, ef := range enabled_features {
				if ef == available_feature.Name {
					available_feature.Enabled = true
					break
				}
			}

			// create feature entry
			err = repo.Set(&types.Feature{
				Name:    available_feature.Name,
				Enabled: available_feature.Enabled,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *App) IsFeatureEnabled(name string) bool {
	f, err := app.Repos.FeatureRepository.GetByName(name)
	if err != nil {
		fmt.Printf("feature query error: %v\n", err)
		return false
	}

	if f != nil {
		return f.Enabled
	} else {
		return false
	}
}
