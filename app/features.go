package app

import (
	"embed"
	"encoding/json"
	"mtui/db"
	"mtui/types"
	"os"
	"strings"
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

func PopulateFeatures(repo *db.FeatureRepository) error {

	available_features, err := GetAvailableFeatures()
	if err != nil {
		return err
	}

	enabled_features := strings.Split(os.Getenv("ENABLE_FEATURES"), ",")

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
