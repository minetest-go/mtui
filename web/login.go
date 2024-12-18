package web

import (
	"encoding/json"
	"fmt"
	"mtui/auth"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OTPCode  string `json:"otp_code"`
}

var tan_map = make(map[string]string)

func (a *Api) TanSetListener(c chan *bridge.CommandResponse) {
	for {
		cmd := <-c
		tc := &command.TanCommand{}
		err := json.Unmarshal(cmd.Data, tc)
		if err != nil {
			fmt.Printf("Tan-listener-error: %s\n", err.Error())
			continue
		}

		if tc.TAN == "" {
			// remove tan
			delete(tan_map, tc.Playername)
		} else {
			// set tan
			tan_map[tc.Playername] = tc.TAN
		}
	}
}

func (a *Api) DoLogout(w http.ResponseWriter, r *http.Request) {
	a.RemoveClaims(w)
}

func (a *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetClaims(r)
	if err == err_unauthorized {
		SendError(w, 401, fmt.Errorf("unauthorized"))
	} else if err != nil {
		SendError(w, 500, err)
	} else if !a.app.MaintenanceMode() && !claims.ApiToken {
		// refresh token
		auth_entry, err := a.app.DBContext.Auth.GetByUsername(claims.Username)
		if err != nil {
			SendError(w, 500, err)
			return
		}
		if auth_entry == nil {
			SendError(w, 404, fmt.Errorf("auth entry not found"))
			return
		}

		claims, err = a.updateToken(w, *auth_entry.ID, claims.Username)
		Send(w, claims, err)

		if claims != nil {
			a.app.CreateUILogEntry(&types.Log{
				Username: claims.Username,
				Event:    "login",
				Message:  fmt.Sprintf("User '%s' refreshed session", claims.Username),
			}, r)
		}

	} else {
		// maintenance mode, send back existing claims
		Send(w, claims, nil)
	}
}

func (a *Api) updateToken(w http.ResponseWriter, id int64, username string) (*types.Claims, error) {
	f, err := a.app.Repos.FeatureRepository.GetByName(types.FEATURE_XBAN)
	if err != nil {
		return nil, fmt.Errorf("could not get feature: %v", err)
	}

	if f.Enabled {
		// consult xban database
		entry, err := a.app.GetOfflineXBanEntry(username)
		if err != nil {
			// just log in error-case, otherwise the login will be blocked
			logrus.WithError(err).Error("could not get xban entry on login")
		}
		if entry != nil && entry.Banned {
			return nil, fmt.Errorf("player is banned (reason: '%s')", entry.Reason)
		}
	}

	privs, err := a.app.DBContext.Privs.GetByID(id)
	if err != nil {
		return nil, err
	}

	priv_arr := make([]string, len(privs))
	for i, p := range privs {
		priv_arr[i] = p.Privilege
	}

	expires := time.Now().Add(7 * 24 * time.Hour)
	claims := &types.Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		Username:   username,
		Privileges: priv_arr,
	}
	return claims, a.SetClaims(w, claims)
}

func (a *Api) DoLogin(w http.ResponseWriter, r *http.Request) {
	req := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	auth_entry, err := a.app.DBContext.Auth.GetByUsername(req.Username)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if auth_entry == nil {
		SendError(w, 404, fmt.Errorf("user not found"))
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: req.Username,
			Event:    "login",
			Message:  fmt.Sprintf("User '%s' tried to login and was not found", req.Username),
		}, r)
		return
	}

	// check password or tan
	tan := tan_map[req.Username]
	if tan == "" {
		// login against the database password

		// legacy first
		legacy_ok := auth.VerifyLegacyPassword(req.Username, req.Password, auth_entry.Password)
		if !legacy_ok {
			// SRP fallback
			salt, verifier, err := auth.ParseDBPassword(auth_entry.Password)
			if err != nil {
				SendError(w, 500, err)
				return
			}

			ok, err := auth.VerifyAuth(req.Username, req.Password, salt, verifier)
			if err != nil {
				SendError(w, 500, err)
				return
			}
			if !ok {
				SendError(w, 401, fmt.Errorf("unauthorized"))
				// create log entry
				a.app.CreateUILogEntry(&types.Log{
					Username: req.Username,
					Event:    "login",
					Message:  fmt.Sprintf("User '%s' provided wrong password", req.Username),
				}, r)
				return
			}
		}
	} else {
		// login with tan
		if tan != req.Password {
			SendError(w, 401, fmt.Errorf("unauthorized"))
			return
		}

		// remove tan (single-use)
		delete(tan_map, req.Username)
	}

	// check otp code if applicable
	privs, err := a.app.DBContext.Privs.GetByID(*auth_entry.ID)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	otp_enabled := false
	for _, priv := range privs {
		if priv.Privilege == "otp_enabled" {
			otp_enabled = true
			break
		}
	}
	if otp_enabled {
		secret_entry, err := a.app.DBContext.ModStorage.Get("otp", []byte(fmt.Sprintf("%s_secret", req.Username)))
		if err != nil {
			SendError(w, 500, err)
			return
		}

		if secret_entry != nil {
			otp_ok, err := totp.ValidateCustom(req.OTPCode, string(secret_entry.Value), time.Now(), totp.ValidateOpts{
				Digits:    6,
				Period:    30,
				Algorithm: otp.AlgorithmSHA1,
			})

			if err != nil {
				SendError(w, 500, err)
				return
			}

			if !otp_ok {
				SendError(w, 403, fmt.Errorf("otp code wrong"))
				// create log entry
				a.app.CreateUILogEntry(&types.Log{
					Username: req.Username,
					Event:    "login",
					Message:  fmt.Sprintf("User '%s' provided wrong otp code", req.Username),
				}, r)
				return
			}
		}
	}

	claims, err := a.updateToken(w, *auth_entry.ID, auth_entry.Name)
	if err != nil {
		SendError(w, 500, fmt.Errorf("token-update failed: %v", err))
		return
	}

	if claims != nil {
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "login",
			Message:  fmt.Sprintf("User '%s' logged in successfully", claims.Username),
		}, r)
	}
	Send(w, claims, err)
}
