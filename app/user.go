package app

import (
	"fmt"
	"mtui/auth"

	dbauth "github.com/minetest-go/mtdb/auth"
)

func (a *App) CreateUser(username, password string, overwrite bool, privs []string) (*dbauth.AuthEntry, error) {
	err := auth.ValidateUsername(username)
	if err != nil {
		return nil, fmt.Errorf("username invalid: %v", err)
	}

	if password == "" {
		return nil, fmt.Errorf("password is empty")
	}

	auth_entry, err := a.DBContext.Auth.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("could not fetch auth entry: %v", err)
	}

	if auth_entry == nil {
		// create new auth entry
		salt, verifier, err := auth.CreateAuth(username, password)
		if err != nil {
			return nil, fmt.Errorf("could not create auth entry: %v", err)
		}

		auth_entry = &dbauth.AuthEntry{
			Name:     username,
			Password: auth.CreateDBPassword(salt, verifier),
		}
		// save to db
		err = a.DBContext.Auth.Create(auth_entry)
		if err != nil {
			return nil, fmt.Errorf("could not save to db: %v", err)
		}
	} else if !overwrite {
		// overwrite not enabled
		return nil, fmt.Errorf("user already exists: '%s'", username)
	}

	existing_privs, err := a.DBContext.Privs.GetByID(*auth_entry.ID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch priv entries: %v", err)
	}

	for _, priv := range privs {
		priv_exists := false
		for _, existing_priv := range existing_privs {
			if existing_priv.Privilege == priv {
				priv_exists = true
				break
			}
		}
		if priv_exists {
			// already there, skip
			continue
		}

		err = a.DBContext.Privs.Create(&dbauth.PrivilegeEntry{
			ID:        *auth_entry.ID,
			Privilege: priv,
		})

		if err != nil {
			return nil, fmt.Errorf("could not create privs: %v", err)
		}
	}

	return auth_entry, nil
}

func (a *App) CreateAdmin(username, password string) (*dbauth.AuthEntry, error) {
	return a.CreateUser(username, password, true, []string{"shout", "server", "interact", "privs", "ban"})
}
