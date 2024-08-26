package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GetRefreshToken() string {
	data := make([]byte, 32)

	if _, err := rand.Read(data); err != nil {
		return ""
	}

	refreshToken := base64.StdEncoding.Strict().EncodeToString(data)

	return refreshToken
}
