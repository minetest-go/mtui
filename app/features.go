package app

import (
	"mtui/db"
	"mtui/types"
)

var feature_list = map[string]bool{
	"mail":          false,
	"areas":         false,
	"skinsdb":       false,
	"luashell":      false,
	"shell":         false,
	"modmanagement": false,
}

func PopulateFeatures(repo *db.FeatureRepository) error {
	for name, enabled := range feature_list {
		feature, err := repo.GetByName(name)
		if err != nil {
			return err
		}
		if feature == nil {
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
