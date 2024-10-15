package db

import (
	"context"
	"errors"
	"math/rand"
	"mtui/types"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OauthAppRepository struct {
	g *gorm.DB
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
	return r.g.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(m).Error
}

func (r *OauthAppRepository) GetAll() ([]*types.OauthApp, error) {
	var list []*types.OauthApp
	err := r.g.Find(&list).Error
	return list, err
}

func (r *OauthAppRepository) GetByName(name string) (*types.OauthApp, error) {
	var list []*types.OauthApp
	err := r.g.Where(&types.OauthApp{Name: name}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *OauthAppRepository) GetByID(id string) (*types.OauthApp, error) {
	var list []*types.OauthApp
	err := r.g.Where(&types.OauthApp{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *OauthAppRepository) Delete(id string) error {
	return r.g.Where(types.OauthApp{ID: id}).Delete(types.OauthApp{}).Error
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
