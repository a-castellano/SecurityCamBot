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

## Allowed actions

For the time being, webcams can be rebooted only.

## Configuration

This bot uses a config file which folder location is defined by environment variable **SECURITY_CAM_BOT_CONFIG_FILE_LOCATION**, inside this folder it must exists a file called **config.toml**.

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

[webcams]
[webcams.cam01]
name= "cam1"
ip = "10.10.10.34"
user = "user"
password = "pass"

[webcams.cam02]
name= "cam2"
ip = "10.10.10.35"
user = "user"
password = "pass"
```

Config files must include the following sections:
### telegram_bot

Defines telegram bot config:
* token -> bot token
* allowed_senders -> list of telegram users allowed to interact with this bot.
  * id -> user ID
  * name -> user name

### Webcams 

Webcams to manage must be set in this section.
* name -> Name for identifying the webcam
* ip -> Webcam IP
* user -> Webcam user
* password -> Webcam password

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
