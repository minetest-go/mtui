package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	repos := db.NewRepositories(setupDB(t))

	mt := &types.MetricType{
		Name: "my_metric",
		Type: "gauge",
		Help: "blah",
	}

	assert.NoError(t, repos.MetricTypeRepository.Insert(mt))

	assert.NoError(t, repos.MetricRepository.Insert(&types.Metric{
		Name:      mt.Name,
		Timestamp: 1000,
		Value:     1,
	}))
	assert.NoError(t, repos.MetricRepository.Insert(&types.Metric{
		Name:      mt.Name,
		Timestamp: 2000,
		Value:     2,
	}))
	assert.NoError(t, repos.MetricRepository.Insert(&types.Metric{
		Name:      mt.Name,
		Timestamp: 3000,
		Value:     3,
	}))
}
