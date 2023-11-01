package connection

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var MinecraftCommands CommandConfiguration

type CommandConfiguration struct {
	Register        string `json:"register"`
	Unregister      string `json:"unregister"`
	ChangePassword  string `json:"changePassword"`
	AddWhitelist    string `json:"addWhitelist"`
	RemoveWhitelist string `json:"removeWhitelist"`
	GetWhitelist    string `json:"getWhitelist"`
	Playtime        string `json:"playtime"`
}

func ImportConfiguration() {
	jsonFile, err := os.Open("commands.json")
	if err != nil {
		log.Panicf("Error opening commands.json: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &MinecraftCommands)
	if err != nil {
		log.Panicf("Error unmarshaling configuration: %v", err)
	}
}
