package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRepo(t *testing.T) {
	_db := setupDB(t)
	repo := db.LogRepository{DB: _db}

	l := &types.Log{}
	assert.NoError(t, repo.Insert(l))
	assert.NotEqual(t, "", l.ID)
	assert.True(t, l.Timestamp > 0)

	from_t := l.Timestamp - 2000
	to_t := l.Timestamp + 2000
	s := &types.LogSearch{
		FromTimestamp: &from_t,
		ToTimestamp:   &to_t,
	}

	count, err := repo.Count(s)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	list, err := repo.Search(s)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, l.ID, list[0].ID)

	l.Category = "xy"
	assert.NoError(t, repo.Update(l))

	list, err = repo.Search(s)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, types.LogCategory("xy"), list[0].Category)
}
