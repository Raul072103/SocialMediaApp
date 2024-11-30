package main

import (
	"SocialMediaApp/internal/db"
	store2 "SocialMediaApp/internal/store"
	"log"
)

func main() {
	database, err := db.New(
		"postgres://admin:adminpassword@localhost:5434/social?sslmode=disable",
		3,
		3,
		"5m")
	defer database.Close()
	if err != nil {
		log.Fatal("Error openning DB for seeding!")
	}

	store := store2.NewStorage(database)

	Seed(db, store)
}
