package main

import (
	"github.com/joho/godotenv"
	"github.com/kitaminka/server-bot/bot"
	"github.com/kitaminka/server-bot/config"
	"github.com/kitaminka/server-bot/connection"
	"log"
	"os"
)

var (
	Token        string
	MongoUri     string
	RconAddress  string
	RconPassword string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}
	Token = os.Getenv("DISCORD_TOKEN")
	MongoUri = os.Getenv("MONGODB_URI")
	RconAddress = os.Getenv("RCON_ADDRESS")
	RconPassword = os.Getenv("RCON_PASSWORD")
}

func main() {
	config.LoadConfig()
	connection.ConnectMongo(MongoUri)
	connection.ConnectRcon(RconAddress, RconPassword)
	connection.GetWhitelist()
	bot.StartBot(Token)
}
