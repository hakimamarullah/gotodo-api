package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("ac815f8cf395c268a004fedec9efbd6788965ce5d237ff5d31b47eaca34918d4280dc52771f3c095de6d17a07071b0270e01334e88679dc57c7bac85ecabfc22")

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
