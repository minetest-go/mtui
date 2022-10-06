package app

import (
	"mtui/auth"
	"time"

	"github.com/minetest-go/mtdb"
	dbauth "github.com/minetest-go/mtdb/auth"
)

func CreateAdminUser(dbctx *mtdb.Context, username, password string) error {
	entry, err := dbctx.Auth.GetByUsername(username)
	if err != nil {
		return err
	}

	if entry != nil {
		// admin user already created
		return nil
	}

	salt, verifier, err := auth.CreateAuth(username, password)
	if err != nil {
		return err
	}

	entry = &dbauth.AuthEntry{
		Name:      username,
		Password:  auth.CreateDBPassword(salt, verifier),
		LastLogin: int(time.Now().Unix()),
	}
	err = dbctx.Auth.Create(entry)
	if err != nil {
		return err
	}

	for _, priv := range []string{"interact", "server"} {
		err = dbctx.Privs.Create(&dbauth.PrivilegeEntry{
			ID:        *entry.ID,
			Privilege: priv,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
