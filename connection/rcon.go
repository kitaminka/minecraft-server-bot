package connection

import (
	"crypto/rand"
	"fmt"
	"github.com/willroberts/minecraft-client"
	"log"
	"math/big"
	"strings"
)

const passwordChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	RconAddress  string
	RconPassword string
	RconClient   *minecraft.Client
)

func ConnectRcon(rconAddress, rconPassword string) {
	RconAddress = rconAddress
	RconPassword = rconPassword
	rconClient, err := connectRconClient(RconAddress, RconPassword)
	if err != nil {
		log.Panicf("Error connecting to RCON: %v", err)
	}
	log.Print("Successfully connected to RCON")
	RconClient = rconClient
}
func ReconnectRcon() {
	log.Print("Trying to reconnect to RCON")
	rconClient, err := connectRconClient(RconAddress, RconPassword)
	if err != nil {
		log.Panicf("Error reconnecting to RCON: %v", err)
	}
	log.Print("Successfully reconnected to RCON")
	RconClient = rconClient
}
func RegisterMinecraftPlayer(minecraftNickname string) (string, error) {
	password := generatePassword()
	err := UnregisterMinecraftPlayer(minecraftNickname)
	if err != nil {
		return password, err
	}
	_, err = sendCommand(fmt.Sprintf("nlogin register %v %v", minecraftNickname, password))
	if err != nil {
		return password, err
	}
	return password, err
}
func UnregisterMinecraftPlayer(minecraftNickname string) error {
	_, err := sendCommand(fmt.Sprintf("nlogin unregister %v", minecraftNickname))
	return err
}
func ResetMinecraftPlayerPassword(minecraftNickname string) (string, error) {
	password := generatePassword()
	err := ChangeMinecraftPlayerPassword(minecraftNickname, password)
	return password, err
}
func ChangeMinecraftPlayerPassword(minecraftNickname, newPassword string) error {
	_, err := sendCommand(fmt.Sprintf("nlogin changepassword %v %v", minecraftNickname, newPassword))
	return err
}
func AddPlayerWhitelist(minecraftNickname string) error {
	message, err := sendCommand(fmt.Sprintf("whitelist add %v", minecraftNickname))
	if err != nil {
		return err
	} else if message.Body == "Player is already whitelisted" {
		err = PlayerAlreadyExistsError
		return err
	}
	return err
}
func RemovePlayerWhitelist(minecraftNickname string) error {
	_, err := sendCommand(fmt.Sprintf("whitelist remove %v", minecraftNickname))
	return err
}
func GetPlayerWhitelist() ([]string, error) {
	message, err := sendCommand("whitelist list")
	if err != nil {
		return nil, err
	} else if len(message.Body) <= 34 {
		return []string{}, err
	}
	playerWhitelist := strings.Split(message.Body[34:], ", ")
	return playerWhitelist, err
}
func GetPlayerPlaytime(minecraftNickname string) (string, error) {
	message, err := sendCommand(fmt.Sprintf("playtime %v", minecraftNickname))
	return message.Body[23+len(minecraftNickname):], err
}
func connectRconClient(rconAddress, rconPassword string) (*minecraft.Client, error) {
	rconClient, err := minecraft.NewClient(rconAddress)
	if err != nil {
		return rconClient, err
	}
	err = rconClient.Authenticate(rconPassword)
	return rconClient, err
}
func sendCommand(command string) (minecraft.Message, error) {
	message, err := RconClient.SendCommand(command)
	if err != nil {
		log.Printf("Error sending command: %v", err)
		ReconnectRcon()
		message, err = RconClient.SendCommand(command)
		return message, err
	}
	return message, nil
}
func generatePassword() string {
	var password string

	for i := 0; i < 8; i++ {
		res, _ := rand.Int(rand.Reader, big.NewInt(61))
		password += string(passwordChars[res.Int64()])
	}

	return password
}
