package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthencticator struct {
	secret   string
	audience string
	issuer   string
}

func NewJWTAuthenticator(secret, audience, issuer string) *JWTAuthencticator {
	return &JWTAuthencticator{secret, audience, issuer}
}

func (a *JWTAuthencticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuthencticator) ValidateToken(token string) (*jwt.Token, error) {

	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing error: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
