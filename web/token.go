package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mtui/types"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TOKEN_COOKIE_NAME = "mtui"

var err_unauthorized = errors.New("unauthorized")

func (api *Api) createCookie(value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    value,
		Path:     api.app.Config.CookiePath,
		Expires:  expires,
		HttpOnly: true,
		Secure:   api.app.Config.CookieSecure,
	}
}

func (api *Api) SetToken(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, api.createCookie(token, expires))
}

func GetToken(r *http.Request) (string, error) {
	// token in cookie
	c, err := r.Cookie(TOKEN_COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		return "", err
	}
	if c != nil {
		return c.Value, nil
	}

	// token in header
	auth_header := r.Header.Get("Authorization")
	if strings.HasPrefix(auth_header, "Bearer ") {
		auth, _ := strings.CutPrefix(auth_header, "Bearer ")
		return auth, nil
	}

	return "", err_unauthorized
}

func (api *Api) RemoveClaims(w http.ResponseWriter) {
	http.SetCookie(w, api.createCookie("", time.Now()))
}

func (api *Api) createToken(claims *types.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(api.app.Config.JWTKey))
}

func (api *Api) SetClaims(w http.ResponseWriter, claims *types.Claims) error {
	token, err := api.createToken(claims)
	if err != nil {
		return err
	}

	api.SetToken(w, token, claims.ExpiresAt.Local())
	return nil
}

func (api *Api) GetClaims(r *http.Request) (*types.Claims, error) {
	var token_string string

	// token in cookie
	c, err := r.Cookie(TOKEN_COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		return nil, fmt.Errorf("error retrieving cookie: %v", err)
	}
	if c != nil {
		token_string = c.Value
	}

	if !api.app.MaintenanceMode() {
		// check api feature
		feature, err := api.app.Repos.FeatureRepository.GetByName(types.FEATURE_API)
		if err != nil {
			return nil, fmt.Errorf("error querying api-feature: %v", err)
		}

		if feature != nil && feature.Enabled {
			// token in header
			auth_header := r.Header.Get("Authorization")
			if strings.HasPrefix(auth_header, "Bearer ") {
				token_string, _ = strings.CutPrefix(auth_header, "Bearer ")
			}
		}
	}

	// no auth found, return 401
	if token_string == "" {
		return nil, err_unauthorized
	}

	// parse string to jwt token
	token, err := jwt.ParseWithClaims(token_string, &types.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(api.app.Config.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err_unauthorized
	}

	claims, ok := token.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("internal error")
	}

	if api.app.IsFeatureEnabled(types.FEATURE_XBAN) {
		// query xban db
		xban, err := api.app.GetOfflineXBanEntry(claims.Username)
		if err != nil {
			return nil, fmt.Errorf("offline xban error: %v", err)
		}

		if xban != nil && xban.Banned {
			return nil, fmt.Errorf("banned player, reason: '%s'", xban.Reason)
		}
	}

	return claims, nil
}

type CreateTokenRequest struct {
	Expiry int64    `json:"expiry"` // millis utc
	Privs  []string `json:"privs"`
}

// Creates an api token for scripting use
// resulting token won't be able to issue other tokens or extend the token lifespan
func (a *Api) CreateToken(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	if claims.ApiToken {
		SendError(w, 403, fmt.Errorf("can't re-issue an api-token with another api-token"))
		return
	}

	ctr := &CreateTokenRequest{}
	err := json.NewDecoder(r.Body).Decode(ctr)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	// check privs
	for _, p := range ctr.Privs {
		if !claims.HasPriv(p) {
			SendError(w, 403, fmt.Errorf("priv not available: '%s'", p))
			return
		}
	}

	t, err := a.createToken(&types.Claims{
		Privileges: ctr.Privs,
		Username:   claims.Username,
		ApiToken:   true,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.UnixMilli(ctr.Expiry)),
		},
	})
	if err != nil {
		SendError(w, 500, fmt.Errorf("create token error: %v", err))
		return
	}

	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "login",
		Message: fmt.Sprintf("User '%s' creates a new token with privs: '%s' expires: %s",
			claims.Username, strings.Join(ctr.Privs, ","), time.UnixMilli(ctr.Expiry).Format(time.RFC3339)),
	}, r)

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(t))
}
