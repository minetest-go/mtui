package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	_, g := setupDB(t)
	repos := db.NewRepositories(g)

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

	// empty list
	from_t := int64(1500)
	to_t := int64(3500)
	name := "bogus"
	s := &types.MetricSearch{
		Name:          &name,
		FromTimestamp: &from_t,
		ToTimestamp:   &to_t,
	}

	count, err := repos.MetricRepository.Count(s)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// list
	s = &types.MetricSearch{
		Name:          &mt.Name,
		FromTimestamp: &from_t,
		ToTimestamp:   &to_t,
	}

	count, err = repos.MetricRepository.Count(s)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// delete
	err = repos.MetricRepository.DeleteBefore(int64(2500))
	assert.NoError(t, err)

	// partial list
	s = &types.MetricSearch{
		Name:          &mt.Name,
		FromTimestamp: &from_t,
		ToTimestamp:   &to_t,
	}

	count, err = repos.MetricRepository.Count(s)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
