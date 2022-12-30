package jobs

import (
	"mtui/app"
	"mtui/types"
	"time"

	"github.com/sirupsen/logrus"
)

func mediaScan(a *app.App) {
	for {
		f, err := a.Repos.FeatureRepository.GetByName(types.FEATURE_MEDIASERVER)
		if err != nil {
			logrus.Errorf("Mediascan getFeature error: %s", err.Error())
		} else if f.Enabled {
			err = a.Mediaserver.ScanDefaultSubdirs(a.WorldDir)
			if err != nil {
				logrus.Errorf("Mediascan scan error: %s", err.Error())
			}
		}
		time.Sleep(time.Minute * 30)
	}
}
