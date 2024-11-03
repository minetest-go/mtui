package app

import (
	"fmt"
	"mtui/db"
	"mtui/mail"
	"mtui/modmanager"
	"path"

	oautherrors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/worldconfig"
	"github.com/sirupsen/logrus"
)

func (a *App) AttachDatabase() error {
	logrus.WithFields(logrus.Fields{
		"worlddir": a.WorldDir,
	}).Info("Attaching database")

	wc, err := worldconfig.Parse(path.Join(a.WorldDir, "world.mt"))
	if err != nil {
		return fmt.Errorf("error parsing world.mt: %v", err)
	}

	supported_backends := map[string]bool{
		worldconfig.BACKEND_POSTGRES: true,
		worldconfig.BACKEND_SQLITE3:  true,
	}

	if !supported_backends[wc[worldconfig.CONFIG_MAP_BACKEND]] {
		fmt.Printf("Warning: unsupported map-backend: '%s' found, defaulting to sqlite!\n", wc[worldconfig.CONFIG_MAP_BACKEND])
		wc[worldconfig.CONFIG_MAP_BACKEND] = worldconfig.BACKEND_SQLITE3
	}

	if !supported_backends[wc[worldconfig.CONFIG_AUTH_BACKEND]] {
		fmt.Printf("Warning: unsupported auth-backend: '%s' found, defaulting to sqlite!\n", wc[worldconfig.CONFIG_AUTH_BACKEND])
		wc[worldconfig.CONFIG_AUTH_BACKEND] = worldconfig.BACKEND_SQLITE3
	}

	if !supported_backends[wc[worldconfig.CONFIG_MOD_STORAGE_BACKEND]] {
		fmt.Printf("Warning: unsupported mod_storage-backend: '%s' found, defaulting to sqlite!\n", wc[worldconfig.CONFIG_MOD_STORAGE_BACKEND])
		wc[worldconfig.CONFIG_MOD_STORAGE_BACKEND] = worldconfig.BACKEND_SQLITE3
	}

	if !supported_backends[wc[worldconfig.CONFIG_PLAYER_BACKEND]] {
		fmt.Printf("Warning: unsupported player-backend: '%s' found, defaulting to sqlite!\n", wc[worldconfig.CONFIG_PLAYER_BACKEND])
		wc[worldconfig.CONFIG_PLAYER_BACKEND] = worldconfig.BACKEND_SQLITE3
	}

	dbctx, err := mtdb.NewWithConfig(a.WorldDir, wc)
	if err != nil {
		return fmt.Errorf("error creating database connection: %v", err)
	}
	a.DBContext = dbctx

	db_, g, err := db.Init(a.WorldDir)
	if err != nil {
		return fmt.Errorf("error initializing local database: %v", err)
	}
	a.DB = db_
	a.G = g

	err = db.Migrate(db_)
	if err != nil {
		return fmt.Errorf("error migrating local database: %v", err)
	}

	a.Repos = db.NewRepositories(g)
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
