package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/meraiku/kode/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestAPIserver(t *testing.T) {
	assert := assert.New(t)

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	t.Run("user creation functional", func(t *testing.T) {

		tests := []struct {
			name      string
			urlPath   string
			wantCode  int
			wantBody  string
			userEmail string
		}{
			{
				name:      "Creating User",
				urlPath:   "/api/users",
				wantCode:  http.StatusCreated,
				wantBody:  "id",
				userEmail: "test@gmail.com",
			},
			{
				name:      "Creating User with existing email",
				urlPath:   "/api/users",
				wantCode:  http.StatusBadRequest,
				userEmail: "test@gmail.com",
			},
			{
				name:     "Creating User without body",
				urlPath:  "/api/users",
				wantCode: http.StatusBadRequest,
			},
			{
				name:      "Creating User with incorrect email",
				urlPath:   "/api/users",
				wantCode:  http.StatusBadRequest,
				userEmail: "1",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				type jsonBody struct {
					Email string `json:"email"`
				}
				b := jsonBody{Email: tt.userEmail}

				code, body := ts.post(t, tt.urlPath, b)

				assert.Equal(tt.wantCode, code)

				if tt.wantBody == "" {
					t.Skip()
				}

				assert.Contains(body, tt.wantBody)

				code, body = ts.get(t, tt.urlPath)

				assert.Equal(http.StatusOK, code)
				assert.Contains(body, tt.userEmail)
			})
		}
	})
	t.Run("token creation", func(t *testing.T) {

		tests := []struct {
			name      string
			urlPath   string
			urlPath2  string
			wantCode  int
			wantBody  string
			userEmail string
			userID    string
		}{
			{
				name:      "From Creation To Tokens",
				urlPath:   "/api/users",
				urlPath2:  "/api/tokens",
				wantCode:  http.StatusCreated,
				wantBody:  "id",
				userEmail: "test1@gmail.com",
			},
			{
				name:     "Create token pair without body",
				urlPath:  "/api/tokens",
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "Create tokens for random id",
				urlPath:  "/api/tokens",
				wantCode: http.StatusBadRequest,
				userID:   uuid.NewString(),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				type jsonBody struct {
					Email string `json:"email"`
					ID    string `json:"id"`
				}

				b := jsonBody{
					Email: tt.userEmail,
					ID:    tt.userID,
				}

				code, body := ts.post(t, tt.urlPath, b)

				assert.Equal(tt.wantCode, code)

				if tt.wantBody == "" {
					t.Skip()
				}

				assert.Contains(body, tt.wantBody)

				type response struct {
					ID           string `json:"id"`
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
				}

				resp := response{}

				json.Unmarshal([]byte(body), &resp)

				err := uuid.Validate(resp.ID)

				assert.Nil(err)

				b.ID = resp.ID

				code, body = ts.post(t, tt.urlPath2, b)

				assert.Equal(code, http.StatusCreated)

				json.Unmarshal([]byte(body), &resp)

				_, err = token.ParseJWT(resp.AccessToken, []byte{})

				assert.Nil(err)
			})
		}
	})
	t.Run("Refresh operation", func(t *testing.T) {

		tests := []struct {
			name      string
			wantCode  int
			userEmail string
			userID    string
		}{
			{
				name:      "From Creation To Refresh Tokens",
				wantCode:  http.StatusCreated,
				userEmail: "testing2@gmail.com",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				type jsonBody struct {
					Email        string `json:"email"`
					ID           string `json:"id"`
					RefreshToken string `json:"refresh_token"`
				}

				type response struct {
					ID           string `json:"id"`
					AccessToken  string `json:"access_token"`
					RefreshToken string `json:"refresh_token"`
				}

				b := jsonBody{Email: tt.userEmail}
				resp := response{}

				code, body := ts.post(t, "/api/users", b)

				assert.Equal(http.StatusCreated, code, body)

				json.Unmarshal([]byte(body), &resp)

				b.ID = resp.ID

				code, body = ts.post(t, "/api/tokens", b)

				assert.Equal(http.StatusCreated, code, fmt.Sprintf("tokens are not created for user with id: %s", b.ID))

				json.Unmarshal([]byte(body), &resp)

				b.RefreshToken = resp.RefreshToken

				code, body = ts.post(t, "/api/tokens/refresh", b)

				assert.Equal(http.StatusCreated, code)

				assert.Contains(body, "refresh", "refresh token is not in the body")

				code, _ = ts.post(t, "/api/tokens/refresh", b)

				assert.Equal(http.StatusBadRequest, code, "refresh tokens are reusale")

				json.Unmarshal([]byte(body), &resp)

				b.RefreshToken = resp.RefreshToken
				b.ID = uuid.NewString()

				code, _ = ts.post(t, "/api/tokens/refresh", b)

				assert.Equal(http.StatusBadRequest, code, "refresh tokens work with random id")
			})
		}
	})

}
