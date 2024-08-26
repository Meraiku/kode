package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = []byte(os.Getenv("JWT_SECRET"))

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetJWT(id string, ip string) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    ip,
		Subject:   id,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(15 * time.Minute)),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := jwtToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseJWT(tokenStr string) (jwt.Claims, error) {

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return key, nil })
	if token != nil {
		if token.Valid {
			return token.Claims, nil
		}
	}

	return nil, err
}
