package main

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/meraiku/kode/internal/database"
	"github.com/meraiku/kode/internal/speller"
	"github.com/meraiku/kode/internal/token"
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

	tokens, err := app.writeTokens(params.Id, cutPort(r.RemoteAddr), w)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			app.respondWithError(w, http.StatusBadRequest, database.ErrNotFound.Error())
			return
		}
		app.errorLog.Print(err)
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

		payload, err := token.ParseJWT(tokens.AccessToken, []byte(app.accessSecret))
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
		if !app.validateRefreshToken(params.RefreshToken, params.Id) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	tokens, err = app.writeTokens(params.Id, cutPort(r.RemoteAddr), w)
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

func (app *application) handleGetNotes(w http.ResponseWriter, r *http.Request) {

	id := app.ctx.Value(key("id"))

	notes, err := app.db.GetUserNotes(id.(string))
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

func (app *application) handlePostNotes(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	params := &parameters{}

	id := app.ctx.Value(key("id"))

	if err := decodeIntoStruct(r, params); err != nil {
		if params.Body == "" || params.Title == "" {
			app.respondWithError(w, http.StatusBadRequest, "missing body")
			return
		}
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	title, err := speller.CheckText(params.Title)
	if err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := speller.CheckText(params.Body)
	if err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	note, err := app.db.CreateNote(id.(string), body, title)
	if err != nil {
		app.errorLog.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params.Body = note.Body
	params.Title = note.Title

	app.respondWithJSON(w, http.StatusCreated, params)
}
