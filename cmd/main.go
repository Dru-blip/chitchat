package main

import (
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewServer(store)
	server.RegisterRoutes()
	server.Start()
}
