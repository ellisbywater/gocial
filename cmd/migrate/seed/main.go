package main

import (
	"log"

	"github.com/ellisbywater/gocial/internal/db"
	"github.com/ellisbywater/gocial/internal/env"
	"github.com/ellisbywater/gocial/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gocial?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
