package main

import (
	"log"
	"time"

	"github.com/ellisbywater/gocial/internal/db"
	"github.com/ellisbywater/gocial/internal/env"
	"github.com/ellisbywater/gocial/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

func main() {
	// Config
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24,
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
