package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ernestosjunior/shortener-url-go/internal/store"
	"github.com/joho/godotenv"

	"github.com/ernestosjunior/shortener-url-go/internal/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
)

func main() {
	initLogger()

	if err := run(); err != nil {
		slog.Error("erro ao executar o serviço", "error", err)
		return
	}

	slog.Info("todos os serviços estão offline")
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	db, err := openDB()
	if err != nil {
		return err
	}

	if err := goose.Up(db, "./internal/store/mysql/migrations"); err != nil {
		slog.Error("erro ao executar migrações no DB", "err", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	store := store.NewStore(rdb)

	server := &http.Server{
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
		Handler:      api.MyHandler(db, store),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
}
