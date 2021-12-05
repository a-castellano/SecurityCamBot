# SecurityCamBot

Telegram bot for managing webcams and receiving alerts from then.

## Install

To Do.

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
* allowed_senders -> list of telegram users allowed to interact with the bot.
  * id -> user ID
  * name -> user name
