package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthRepo(t *testing.T) {
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

	// create repo
	repo := NewAuthRepository(auth_db)
	assert.NotNil(t, repo)

	// non-existing entry
	entry, err := repo.GetByUsername("bogus")
	assert.NoError(t, err)
	assert.Nil(t, entry)

	// create entry
	new_entry := &AuthEntry{
		Name:      "createduser",
		Password:  "blah",
		LastLogin: 456,
	}
	assert.NoError(t, repo.Create(new_entry))
	assert.NotNil(t, new_entry.ID)

	// check newly created entry
	entry, err = repo.GetByUsername("createduser")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, new_entry.Name, entry.Name)
	assert.Equal(t, new_entry.Password, entry.Password)
	assert.Equal(t, *new_entry.ID, *entry.ID)
	assert.Equal(t, new_entry.LastLogin, entry.LastLogin)

	// change things
	new_entry.Name = "x"
	new_entry.Password = "y"
	new_entry.LastLogin = 123
	assert.NoError(t, repo.Update(new_entry))
	entry, err = repo.GetByUsername("x")
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, new_entry.Name, entry.Name)
	assert.Equal(t, new_entry.Password, entry.Password)
	assert.Equal(t, *new_entry.ID, *entry.ID)
	assert.Equal(t, new_entry.LastLogin, entry.LastLogin)

	// remove new user
	assert.NoError(t, repo.Delete(*new_entry.ID))
	entry, err = repo.GetByUsername("x")
	assert.NoError(t, err)
	assert.Nil(t, entry)

}
