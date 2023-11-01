# minecraft-server-bot
**Discord bot** for managing a whitelist and authorization system on Minecraft server.
### Installation
- Clone the repository:
```bash
git clone https://github.com/kitaminka/minecraft-server-bot.git
```
- Create **.env** file and fill it. Example of .env file you can see in **.env.example** file.
- By default, bot works with **[AuthMeReloaded](https://github.com/AuthMe/AuthMeReloaded/)** plugin, but you can change **config.json** to work with any other authorization plugin.
- Start the bot:
```bash
go run main.go
```
- Add bot to your Discord server.
- Using **/send-whitelist** command create whitelist info message.
- Using **/settings set** command set other bot settings.