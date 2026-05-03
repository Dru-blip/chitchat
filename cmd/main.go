package main

import (
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {
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
	defer store.Db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	stmp_mailer, err := mailer.New()
	if err != nil {
		log.Fatal(err)
	}
	server, err := api.NewServer(store, stmp_mailer, rdb)
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterRoutes()
	server.Start()
}
