package connection

import (
	"github.com/bwmarrin/discordgo"
)

func RegisterPlayer(member *discordgo.Member, minecraftNickname string) (string, error) {
	err := CreatePlayer(member, minecraftNickname)
	if err != nil {
		return "", err
	}

	err = AddPlayerWhitelist(minecraftNickname)
	if err != nil {
		return "", err
	}

	password, err := RegisterMinecraftPlayer(minecraftNickname)
	return password, err
}
func UnregisterPlayer(member *discordgo.Member) (Player, error) {
	player, err := GetPlayerByDiscord(member)
	if err != nil {
		return player, err
	}

	_ = UnregisterMinecraftPlayer(player.MinecraftNickname)
	_ = RemovePlayerWhitelist(player.MinecraftNickname)
	_ = DeletePlayer(member)

	return player, err
}
