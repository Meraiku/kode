package token

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTGeneration(t *testing.T) {

	assert := assert.New(t)

	idWant := "1"
	ipWant := "ip"

	got, err := GetJWT(idWant, ipWant)

	assert.Nil(err)

	assert.NotEqual(got, "")

	payload, err := ParseJWT(got)

	assert.Nil(err)

	id, _ := payload.GetSubject()
	ip, _ := payload.GetIssuer()

	assert.Equal(idWant, id)
	assert.Equal(ipWant, ip)
}

func TestRefreshTokenGeneration(t *testing.T) {

	assert := assert.New(t)

	refresh := GetRefreshToken()

	assert.NotEmpty(refresh)

	decoded, err := base64.StdEncoding.DecodeString(refresh)

	assert.Nil(err)
	assert.NotEmpty(decoded)
}
