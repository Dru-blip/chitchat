package main

import (
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	//TODO: should remove this from init function
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterRoutes()
	server.Start()
}
