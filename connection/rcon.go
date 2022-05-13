package connection

import (
	"fmt"
	"github.com/willroberts/minecraft-client"
	"log"
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

	RconClient = rconClient
}

func GetWhitelist() {
	message, _ := RconClient.SendCommand("whitelist list")
	fmt.Println(message.Body)
}
