package web

import (
	"encoding/base64"
	"net/http"
)

func (a *Api) OAuthAuthHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := a.GetClaims(r)
	if err != nil {
		return "", err
	}
	return c.Username, nil
}

func (a *Api) OauthUserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	c, err := a.GetClaims(r)
	if err != nil {
		r.ParseForm()
		returnto := base64.URLEncoding.EncodeToString([]byte(r.URL.String()))
		w.Header().Set("Location", "./#/login?return_to="+returnto)
		w.WriteHeader(http.StatusFound)
		return
	}

	return c.Username, nil
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
