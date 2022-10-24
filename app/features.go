package app

import (
	"mtui/db"
	"mtui/types"
	"os"
	"strings"
)

var feature_list = map[string]bool{
	"mail":          false,
	"areas":         false,
	"skinsdb":       false,
	"luashell":      false,
	"shell":         false,
	"modmanagement": false,
	"mediaserver":   false,
}

func PopulateFeatures(repo *db.FeatureRepository) error {
	enabled_features := strings.Split(os.Getenv("ENABLE_FEATURES"), ",")

	for name, enabled := range feature_list {
		feature, err := repo.GetByName(name)
		if err != nil {
			return err
		}
		if feature == nil {
			// check if the feature was enabled in the env-vars
			for _, ef := range enabled_features {
				if ef == name {
					enabled = true
					break
				}
			}

			// create feature entry
			err = repo.Set(&types.Feature{
				Name:    name,
				Enabled: enabled,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
