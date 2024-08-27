package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/meraiku/kode/internal/token"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, errors.New("DB_URL environment varible is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("DB Ping failed")
		return nil, err
	}

	return db, nil
}

func connectRedis(ctx context.Context) (*redis.Client, error) {

	opt, _ := redis.ParseURL(os.Getenv("RDB_URL"))

	rdb := redis.NewClient(opt)

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func (app *application) respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		ErrMsg string `json:"error"`
	}

	msgToReturn := returnError{
		ErrMsg: msg,
	}

	app.respondWithJSON(w, code, msgToReturn)
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		app.errorLog.Printf("Error marshal json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func (app *application) writeTokens(id, ip string, w http.ResponseWriter) (*token.Tokens, error) {

	access, err := token.GetJWT(id, ip)
	if err != nil {
		return nil, err
	}

	refresh := token.GetRefreshToken()

	tokens := &token.Tokens{AccessToken: access, RefreshToken: refresh}

	cryptoRefresh, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := app.db.UpdateUserInfo(id, string(cryptoRefresh), app.ctx); err != nil {
		return nil, err
	}

	if err := app.cache.SetTokens(id, tokens, app.ctx); err != nil {
		return nil, err
	}

	tokensToCookies(w, tokens, id)

	return tokens, nil
}

func (app *application) validateRefreshToken(refreshToken, id string) bool {
	user, err := app.db.GetUserByID(id, app.ctx)
	if err != nil {
		return false
	}

	if user.RefreshToken == nil {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.RefreshToken), []byte(refreshToken)); err != nil {
		return false
	}

	return true
}

func tokensToCookies(w http.ResponseWriter, tokens *token.Tokens, id string) {
	cookieAccess := &http.Cookie{
		Name:   "access",
		Value:  tokens.AccessToken,
		MaxAge: 15 * 60,
	}
	cookieRefresh := &http.Cookie{
		Name:   "refresh",
		Value:  tokens.RefreshToken,
		MaxAge: 60 * 60 * 24,
	}
	cookieUserId := &http.Cookie{
		Name:   "id",
		Value:  id,
		MaxAge: 60 * 60 * 24,
	}

	setCookies(w, []http.Cookie{*cookieAccess, *cookieRefresh, *cookieUserId})
}

func setCookies(w http.ResponseWriter, cookies []http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, &cookie)
	}
}

func (app *application) sendEmailNotification(recipient string) error {
	server := os.Getenv("SMTP_SERVER")

	if server == "" {
		app.errorLog.Println("continuing without sending email. Add SMTP_SERVER to .env")
		return nil
	}

	from := os.Getenv("SMTP_NAME")
	pass := os.Getenv("SMTP_PASS")

	msg := "From: " + from + "\n" +
		"To: " + recipient + "\n" +
		"Subject: New location\n\n" +
		"Trying access account from new location"

	auth := smtp.PlainAuth("", from, pass, server)

	err := smtp.SendMail(server+":587", auth, from, []string{recipient}, []byte(msg))
	if err != nil {
		return err
	}
	app.infoLog.Printf("email sent to %s\n", recipient)

	return nil
}

func decodeIntoStruct(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding parameters: %s", err)
	}

	return nil
}

func cutPort(ip string) string {
	addr := strings.Split(ip, ":")

	return addr[0]
}
