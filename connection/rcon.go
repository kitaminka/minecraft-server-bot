package connection

import (
	"github.com/willroberts/minecraft-client"
	"log"
	"strings"
)

var RconClient *minecraft.Client

func ConnectRcon(rconAddress, rconPassword string) {
	rconClient, err := minecraft.NewClient(rconAddress)
	if err != nil {
		log.Panicf("Error connecting to RCON: %v", err)
	}

	err = rconClient.Authenticate(rconPassword)
	if err != nil {
		log.Panicf("Error RCON authenticating: %v", err)
	}
	log.Print("Successfully connected to RCON")

	RconClient = rconClient
}
func GetPlayerWhitelist() []string {
	message, err := RconClient.SendCommand("whitelist list")
	if err != nil {
		log.Printf("Error sending command: %v", err)
		return nil
	}
	players := strings.Split(message.Body[33:], ", ")
	return players
}
