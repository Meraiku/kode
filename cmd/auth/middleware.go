package main

import (
	"net/http"
	"time"

	"github.com/meraiku/kode/internal/token"
)

type authHandler func(http.ResponseWriter, *http.Request, string)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticateUser(next authHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access")
		if err != nil {
			refreshToken, err := r.Cookie("refresh")
			if err != nil {
				app.respondWithError(w, http.StatusUnauthorized, "missing Api Keys in cookies")
				return
			}
			id, err := r.Cookie("id")
			if err != nil {
				app.respondWithError(w, http.StatusUnauthorized, "missing Api Keys in cookies")
				return
			}

			if app.validateRefreshToken(refreshToken.Value, id.Value) {
				tokens, _ := app.writeTokens(id.Value, cutPort(r.RemoteAddr), w)
				cookie.Value = tokens.AccessToken

			} else {
				app.respondWithError(w, http.StatusUnauthorized, "missing Api Keys in cookies")
				return
			}
		}

		payload, err := token.ParseJWT(cookie.Value)
		if err != nil {
			app.respondWithError(w, http.StatusUnauthorized, "error parsing Api Key")
			return
		}

		id, err := payload.GetSubject()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expirationTime, err := payload.GetExpirationTime()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expired := expirationTime.Sub(time.Now().UTC())
		switch {
		case expired < time.Second:
			app.respondWithError(w, http.StatusUnauthorized, "Api Key expired")
			return
		case expired < 5*time.Minute:
			_, err := app.writeTokens(id, cutPort(r.RemoteAddr), w)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		next(w, r, id)
	})
}
