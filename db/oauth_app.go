package db

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"mtui/types"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type OauthAppRepository struct {
	dbu *dbutil.DBUtil[*types.OauthApp]
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (r *OauthAppRepository) Set(m *types.OauthApp) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	if m.Secret == "" {
		m.Secret = randSeq(16)
	}
	if m.Created == 0 {
		m.Created = time.Now().Unix()
	}
	return r.dbu.InsertOrReplace(m)
}

func (r *OauthAppRepository) GetAll() ([]*types.OauthApp, error) {
	return r.dbu.SelectMulti("")
}

func (r *OauthAppRepository) GetByName(name string) (*types.OauthApp, error) {
	f, err := r.dbu.Select("where name = %s", name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return f, err
	}
}

func (r *OauthAppRepository) GetByID(id string) (*types.OauthApp, error) {
	f, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return f, err
	}
}

func (r *OauthAppRepository) Delete(id string) error {
	return r.dbu.Delete("where id = %s", id)
}

type OAuthAppStore struct {
	Repo *OauthAppRepository
}

func (s *OAuthAppStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	// id == name
	app, err := s.Repo.GetByName(id)
	if err != nil {
		return nil, err
	} else if !app.Enabled {
		return nil, errors.New("app not enabled")
	} else {
		return app, nil
	}
}
