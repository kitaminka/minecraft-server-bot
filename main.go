package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/bot"
	"github.com/kitaminka/server-bot/database"
	"log"
	"os"
)

var (
	Token             string
	MongoUri          string
	MongoDatabaseName string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGODB_URI")
	MongoDatabaseName = os.Getenv("MONGODB_DATABASE")
}

func main() {
	database.Connect(MongoUri, MongoDatabaseName)
	bot.StartBot(Token)
}
