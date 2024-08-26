package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/meraiku/kode/internal/cache"
	"github.com/meraiku/kode/internal/database"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       database.Store
	cache    cache.Cache
	ctx      context.Context
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		errorLog.Print("Error loading .env file!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	addr := ":" + port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB()
	if err != nil {
		errorLog.Fatalf("Error connecting DB: %s", err)
	}
	defer db.Close()

	rdb, err := connectRedis(ctx)
	if err != nil {
		errorLog.Fatalf("Error connecting Redis: %s", err)
	}
	defer rdb.Close()

	cfg := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       database.NewDB(db),
		ctx:      ctx,
		cache:    cache.NewCache(rdb),
	}

	srv := http.Server{
		ErrorLog:     errorLog,
		Addr:         addr,
		Handler:      cfg.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server running at port %s\n", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
