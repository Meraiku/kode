package token

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTGeneration(t *testing.T) {

	tok := NewTokens()
	id := uuid.NewString()
	secret := "aya"

	token, err := tok.ID(id).ExpiredAt(time.Second).Generate([]byte(secret))

	tok.Claims, err = ParseJWT(token, []byte(secret))

	assert.Nil(t, err, "expected valid token, got invalid")

	assert.Equal(t, id, tok.Claims.ID, "expected %s, got %s", id, tok.Claims.ID)

	time.Sleep(2 * time.Second)

	_, err = ParseJWT(token, []byte(secret))

	assert.NotNil(t, err, "expected invalid token, got valid")
}

func TestRefreshTokenGeneration(t *testing.T) {

	assert := assert.New(t)

	refresh := GetRefreshToken()

	assert.NotEmpty(refresh)

	decoded, err := base64.StdEncoding.DecodeString(refresh)

	assert.Nil(err)
	assert.NotEmpty(decoded)
}
