// Package authclaim hosts the JWT claim type in a neutral package so both
// middleware and service layers can depend on it without service importing
// middleware (which created an inverted dependency).
package authclaim

import "github.com/golang-jwt/jwt/v5"

// Claims is the shared JWT claim structure.
type Claims struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
