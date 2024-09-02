package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	Claims       *Claims `json:"-"`
}

type Claims struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func NewTokens() *Tokens {
	return &Tokens{
		Claims: &Claims{},
	}
}

func (t *Tokens) Generate(secret []byte) (string, error) {

	t.Claims.UID = uuid.NewString()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, t.Claims)

	token, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}

	t.Claims = &Claims{}

	return token, nil
}

func (t *Tokens) ID(id string) *Tokens {
	t.Claims.ID = id

	return t
}

func (t *Tokens) Issuer(issuer string) *Tokens {
	t.Claims.Issuer = issuer

	return t
}

func (t *Tokens) ExpiredAt(expirationTime time.Duration) *Tokens {
	t.Claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expirationTime).UTC())
	return t
}

func ParseJWT(tokenStr string, secret []byte) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return secret, nil })
	if token != nil {
		if token.Valid {
			return claims, nil
		}
	}

	return nil, err
}
