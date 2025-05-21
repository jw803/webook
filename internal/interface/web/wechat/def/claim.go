package def

import "github.com/golang-jwt/jwt/v5"

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
