package main

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/meraiku/kode/internal/database"
	"github.com/meraiku/kode/internal/token"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) handleGetTokens(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Id string `json:"id"`
	}

	params := parameters{}

	if err := decodeIntoStruct(r, &params); err != nil {
		if params.Id == "" {
			app.respondWithError(w, http.StatusBadRequest, "request body is empty")
			return
		}
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokens, err := app.writeTokens(params.Id, cutPort(r.RemoteAddr))
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.respondWithError(w, http.StatusBadRequest, database.ErrNotFound.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.respondWithJSON(w, http.StatusCreated, &tokens)
}

func (app *application) handleRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Id           string `json:"id"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}

	if err := decodeIntoStruct(r, &params); err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokens, err := app.cache.GetTokens(params.Id, app.ctx)
	switch err {
	case nil:
		if tokens.RefreshToken != params.RefreshToken {
			app.respondWithError(w, http.StatusBadRequest, "invalid operation")
			return
		}

		payload, err := token.ParseJWT(tokens.AccessToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ip, err := payload.GetIssuer()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		idJWT, err := payload.GetSubject()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if idJWT != params.Id {
			app.respondWithError(w, http.StatusBadRequest, "Trying refresh tokens for another user")
			app.errorLog.Printf("Trying refresh token for id=%s with id=%s\n", idJWT, params.Id)
			return
		}

		if ip != cutPort(r.RemoteAddr) {
			user, err := app.db.GetUserByID(params.Id, app.ctx)
			if err != nil {
				if errors.Is(err, database.ErrNotFound) {
					app.respondWithError(w, http.StatusBadRequest, database.ErrNotFound.Error())
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := app.sendEmailNotification(user.Email); err != nil {
				app.errorLog.Println(err)
			}
			return
		}

	default:
		user, err := app.db.GetUserByID(params.Id, app.ctx)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				app.respondWithError(w, http.StatusBadRequest, database.ErrNotFound.Error())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if user.RefreshToken == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*user.RefreshToken), []byte(params.RefreshToken)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	tokens, err = app.writeTokens(params.Id, cutPort(r.RemoteAddr))
	if err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.respondWithJSON(w, http.StatusCreated, &tokens)
}

func (app *application) handlePostUsers(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	params := parameters{}

	if err := decodeIntoStruct(r, &params); err != nil {
		if params.Email == "" {
			app.respondWithError(w, http.StatusBadRequest, "request body is empty")
			return
		}
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !reg.Match([]byte(params.Email)) {
		app.respondWithError(w, http.StatusBadRequest, "incorrect email form")
		return
	}

	id, err := app.db.CreateUser(params.Email, app.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			app.respondWithError(w, http.StatusBadRequest, "email already taken")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type out struct {
		Id string `json:"id"`
	}

	app.respondWithJSON(w, http.StatusCreated, out{Id: id})
}

func (app *application) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.db.GetUsers(app.ctx)
	if err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.respondWithJSON(w, http.StatusOK, users)
}

func (app *application) handleGetNotes(w http.ResponseWriter, r *http.Request, id string) {

	notes, err := app.db.GetUserNotes(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.respondWithError(w, http.StatusBadRequest, "notes not found")
			return
		}
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.respondWithJSON(w, http.StatusOK, notes)
}

func (app *application) handlePostNotes(w http.ResponseWriter, r *http.Request, id string) {
	type parameters struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	params := &parameters{}

	if err := decodeIntoStruct(r, params); err != nil {
		if params.Body == "" || params.Title == "" {
			app.respondWithError(w, http.StatusBadRequest, "missing body")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := app.db.CreateNote(id, params.Body, params.Title); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.respondWithJSON(w, http.StatusCreated, nil)
}
