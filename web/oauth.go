package web

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a *Api) OAuthAuthHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := a.GetClaims(r)
	if err != nil {
		return "", err
	}
	return c.Username, nil
}

func (a *Api) OauthUserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	logrus.WithFields(logrus.Fields{
		"Request": r,
	}).Info("OauthUserAuthorizationHandler")

	c, err := a.GetClaims(r)
	if err != nil {
		return "", err
	}

	auth, err := a.app.DBContext.Auth.GetByUsername(c.Username)
	if err != nil {
		return "", err
	}

	if auth == nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	return auth.Name, nil
}

func (a *Api) OAuthAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	err := a.app.OAuthServer.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *Api) OAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	a.app.OAuthServer.HandleTokenRequest(w, r)
}
