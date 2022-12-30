package app

import (
	"mtui/db"
	"mtui/types"
	"os"
	"strings"
)

var feature_list = map[string]bool{
	types.FEATURE_MAIL:          false,
	types.FEATURE_AREAS:         false,
	types.FEATURE_SKINSDB:       false,
	types.FEATURE_LUASHELL:      false,
	types.FEATURE_SHELL:         false,
	types.FEATURE_MODMANAGEMENT: false,
	types.FEATURE_MEDIASERVER:   false,
	types.FEATURE_XBAN:          false,
	types.FEATURE_MONITORING:    false,
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
