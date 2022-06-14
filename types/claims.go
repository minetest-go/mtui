package types

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	*jwt.RegisteredClaims
	Username   string   `json:"username"`
	Privileges []string `json:"privileges"`
}

func (c *Claims) HasPriv(priv string) bool {
	for _, p := range c.Privileges {
		if p == priv {
			return true
		}
	}
	return false
}
