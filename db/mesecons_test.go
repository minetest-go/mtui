package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMesecons(t *testing.T) {
	repos := db.NewRepositories(setupDB(t))

	m1 := &types.Mesecons{
		PosKey:     "1,2,3",
		X:          1,
		Y:          2,
		Z:          3,
		PlayerName: "singleplayer",
	}
	assert.NoError(t, repos.MeseconsRepo.Save(m1))

	m2 := &types.Mesecons{
		PosKey:     "1,2,4",
		X:          1,
		Y:          2,
		Z:          4,
		PlayerName: "otherplayer",
	}
	assert.NoError(t, repos.MeseconsRepo.Save(m2))

	list, err := repos.MeseconsRepo.GetByPlayerName(m1.PlayerName)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	m, err := repos.MeseconsRepo.GetByPoskey(m2.PosKey)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}
