package connection

func RegisterPlayer(userId string, minecraftNickname string) (string, error) {
	err := CreatePlayer(userId, minecraftNickname)
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
func UnregisterPlayer(userId string) (Player, error) {
	player, err := GetPlayerByDiscord(userId)
	if err != nil {
		return player, err
	}

	_ = UnregisterMinecraftPlayer(player.MinecraftNickname)
	_ = RemovePlayerWhitelist(player.MinecraftNickname)
	_ = DeletePlayer(userId)

	return player, err
}
