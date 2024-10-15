package app

import (
	"fmt"
	"mtui/db"
	"mtui/mail"
	"mtui/modmanager"

	oautherrors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

func (a *App) AttachDatabase() error {
	logrus.WithFields(logrus.Fields{
		"worlddir": a.WorldDir,
	}).Info("Attaching database")

	dbctx, err := mtdb.New(a.WorldDir)
	if err != nil {
		return err
	}
	a.DBContext = dbctx

	db_, _, err := db.Init(a.WorldDir)
	if err != nil {
		return err
	}
	a.DB = db_

	err = db.Migrate(db_)
	if err != nil {
		return err
	}

	a.Repos = db.NewRepositories(db_)
	a.Mail = mail.New(dbctx)
	a.ModManager = modmanager.New(a.WorldDir, a.Repos.ModRepo)

	// oauth setup

	oauth_mgr := manage.NewDefaultManager()
	oauth_mgr.MustTokenStorage(store.NewMemoryTokenStore())
	oauth_mgr.MapClientStorage(&db.OAuthAppStore{Repo: a.Repos.OauthAppRepo})
	oauth_mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(a.Config.JWTKey), jwt.SigningMethodHS512))
	a.OAuthMgr = oauth_mgr

	oauth_srv := server.NewDefaultServer(oauth_mgr)
	oauth_srv.SetInternalErrorHandler(func(err error) (re *oautherrors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Internal error")
		return
	})
	oauth_srv.SetResponseErrorHandler(func(re *oautherrors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Response error")
	})
	a.OAuthServer = oauth_srv

	return nil
}

func (a *App) DetachDatabase() error {
	logrus.WithFields(logrus.Fields{
		"worlddir": a.WorldDir,
	}).Info("Detaching database")

	a.OAuthServer = nil
	a.OAuthMgr = nil
	a.ModManager = nil
	a.Mail = nil
	a.Repos = nil
	err := a.DB.Close()
	if err != nil {
		return fmt.Errorf("could not close database: %v", err)
	}
	a.DBContext.Close()
	return nil
}
