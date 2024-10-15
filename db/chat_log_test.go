package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatLogRepo(t *testing.T) {
	_db, g := setupDB(t)
	repo := db.NewRepositories(_db, g).ChatLogRepo

	assert.NoError(t, repo.Insert(&types.ChatLog{Timestamp: 100, Channel: "main", Name: "player1", Message: "msg1"}))
	assert.NoError(t, repo.Insert(&types.ChatLog{Timestamp: 200, Channel: "main", Name: "player2", Message: "msg2"}))
	assert.NoError(t, repo.Insert(&types.ChatLog{Timestamp: 150, Channel: "other_chan", Name: "player1", Message: "msg1"}))

	list, err := repo.GetLatest("main", 100)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "msg1", list[0].Message)
	assert.Equal(t, "player1", list[0].Name)
	assert.Equal(t, "msg2", list[1].Message)
	assert.Equal(t, "player2", list[1].Name)

	list, err = repo.Search("main", 99, 199)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "msg1", list[0].Message)
	assert.Equal(t, "player1", list[0].Name)

	assert.NoError(t, repo.DeleteBefore(199))

	list, err = repo.GetLatest("main", 100)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "msg2", list[0].Message)
	assert.Equal(t, "player2", list[0].Name)
}
