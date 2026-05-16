package main

import (
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"chitchat/internal/mqttclient"
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
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer rdb.Close()

	stmp_mailer, err := mailer.New()
	if err != nil {
		log.Fatal(err)
	}

	mqttClient, err := mqttclient.New()
	if err != nil {
		log.Fatal(err)
	}

	server, err := api.NewApp(store, stmp_mailer, rdb, mqttClient)
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterRoutes()
	server.Start()
}
