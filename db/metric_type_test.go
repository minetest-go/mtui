package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricType(t *testing.T) {
	repos := db.NewRepositories(setupDB(t))

	mt := &types.MetricType{
		Name: "my_metric",
		Type: "gauge",
		Help: "blah",
	}

	// create
	assert.NoError(t, repos.MetricTypeRepository.Insert(mt))

	// read
	mt2, err := repos.MetricTypeRepository.GetByName("junk")
	assert.NoError(t, err)
	assert.Nil(t, mt2)

	// read
	mt2, err = repos.MetricTypeRepository.GetByName("my_metric")
	assert.NoError(t, err)
	assert.NotNil(t, mt2)
	assert.Equal(t, mt.Name, mt2.Name)
	assert.Equal(t, mt.Type, mt2.Type)
	assert.Equal(t, mt.Help, mt2.Help)

	// update
	mt.Type = "counter"
	assert.NoError(t, repos.MetricTypeRepository.Insert(mt))

	// read
	mt2, err = repos.MetricTypeRepository.GetByName("my_metric")
	assert.NoError(t, err)
	assert.NotNil(t, mt2)
	assert.Equal(t, mt.Name, mt2.Name)
	assert.Equal(t, mt.Type, mt2.Type)
	assert.Equal(t, mt.Help, mt2.Help)

	// read all
	list, err := repos.MetricTypeRepository.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
}
