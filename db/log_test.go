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

	l := &types.Log{
		Event:    "myevent",
		Category: types.CategoryMinetest,
	}
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
	assert.Equal(t, 1, count)

	list, err := repo.Search(s)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, l.ID, list[0].ID)

	l.Category = types.CategoryUI
	assert.NoError(t, repo.Update(l))

	list, err = repo.Search(s)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, types.CategoryUI, list[0].Category)

	events, err := repo.GetEvents(types.CategoryUI)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, "myevent", events[0])

	events, err = repo.GetEvents(types.CategoryMinetest)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Equal(t, 0, len(events))
}
