# SecurityCamBot

Telegram bot for managing webcams and receiving alerts from then.

## Install

Add Widmaker repo ans install **windmaker-security-cam-bot**:
```bash
wget -O - https://packages.windmaker.net/WINDMAKER-GPG-KEY.pub | sudo apt-key add -
sudo add-apt-repository "deb http://packages.windmaker.net/ focal main"
sudo apt-get update
sudo apt-get install windmaker-security-cam-bot
```

## Configuration

This bot uses a config file which folder location is defined by environment variable *SECURITY_CAM_BOT_CONFIG_FILE_LOCATION*, inside this folder it must exists a file called *config.toml*.

```toml
[telegram_bot]
token = "token"

[telegram_bot.allowed_senders]
[telegram_bot.allowed_senders.alice]
name = "Alice"
id = 12

[telegram_bot.allowed_senders.bob]
name = "Bob"
id = 13
```

Config files must include the following sections:
### telegram_bot

Defines telegram bot config:
* token -> bot token
* allowed_senders -> list of telegram users allowed to interact with this bot.
  * id -> user ID
  * name -> user name

## Systemd service setup

After saving config in **/etc/windmaker-security-cam-bot/config.toml** systemd service can be enabled:
```bash
sudo /bin/systemctl daemon-reload
sudo /bin/systemctl enable windmaker-security-cam-bot
sudo /bin/systemctl start windmaker-security-cam-bot
```

## Logging

This bot wites logs to syslog:
```
Dec  5 15:51:29 metatron security-cam-bot[8220]: Blocked message received from sender 112897183.
Dec  5 15:56:40 metatron security-cam-bot[8962]: /hello received from sender Bob.
```
