package db_test

import (
	"mtui/db"
	"mtui/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOauthRepo(t *testing.T) {
	DB := setupDB(t)
	repo := db.OauthAppRepository{DB: DB}

	// create
	err := repo.Set(&types.OauthApp{
		Name:         "my-app",
		Enabled:      true,
		Created:      time.Now().Unix(),
		RedirectURLS: "",
		Secret:       "",
		AllowPrivs:   "",
	})
	assert.NoError(t, err)

	// get by name
	app, err := repo.GetByName("my-app")
	assert.NoError(t, err)
	assert.NotNil(t, app)

	// get by id
	app, err = repo.GetByID(app.ID)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	// delete
	err = repo.Delete(app.ID)
	assert.NoError(t, err)

	// get nil
	app, err = repo.GetByName("my-app")
	assert.NoError(t, err)
	assert.Nil(t, app)
}
