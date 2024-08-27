package main

import (
	"net/http"
	"strings"

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

		header, ok := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if !ok {
			app.respondWithError(w, http.StatusUnauthorized, "missing Api Key in headers")
			return
		}

		payload, err := token.ParseJWT(header)
		if err != nil {
			app.respondWithError(w, http.StatusUnauthorized, "error parsing Api Key")
			return
		}

		id, err := payload.GetSubject()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next(w, r, id)
	})
}
