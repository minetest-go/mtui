package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivRepo(t *testing.T) {
	// temp files
	dbfile, err := os.CreateTemp(os.TempDir(), "auth.sqlite")
	assert.NoError(t, err)
	assert.NotNil(t, dbfile)

	// open db
	auth_db, err := sql.Open("sqlite", dbfile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, auth_db)

	// migrate
	err = MigrateAuth(auth_db)
	assert.NoError(t, err)

	// open db
	repo := NewUserPrivilegeRepository(auth_db)
	assert.NotNil(t, repo)

	// create
	assert.NoError(t, repo.Create(&PrivilegeEntry{ID: 2, Privilege: "stuff"}))

	// verify
	list, err := repo.GetByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	privs := make(map[string]bool)
	for _, e := range list {
		privs[e.Privilege] = true
	}
	assert.True(t, privs["stuff"])

	// delete
	assert.NoError(t, repo.Delete(2, "stuff"))

	list, err = repo.GetByID(2)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 0, len(list))

}
