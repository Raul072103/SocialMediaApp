package main

import (
	"SocialMediaApp/internal/db"
	"SocialMediaApp/internal/env"
	"SocialMediaApp/internal/store"
	"SocialMediaApp/scripts"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5434/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	repo := store.NewStorage(conn)

	scripts.Seed(repo)
}
