package types

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	*jwt.RegisteredClaims
	Username   string   `json:"username"`
	Privileges []string `json:"privileges"`
	ApiToken   bool     `json:"api_token"`
}

func (c *Claims) HasPriv(priv string) bool {
	for _, p := range c.Privileges {
		if p == priv {
			return true
		}
	}
	return false
}
