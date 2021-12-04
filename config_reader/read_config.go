package config

import (
	"errors"

	viperLib "github.com/spf13/viper"
)

type TelegramBot struct {
	Token string
}

type Config struct {
	TelegramBot TelegramBot
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "SECURITY_CAM_BOT_CONFIG_FILE_LOCATION"

	telegramBotVariables := []string{"token"}

	viper := viperLib.New()

	//Look for config file location defined as env var
	viper.BindEnv(envVariable)
	configFileLocation = viper.GetString(envVariable)
	if configFileLocation == "" {
		// Get config file from default location
		return config, errors.New(errors.New("Environment variable SECURITY_CAM_BOT_CONFIG_FILE_LOCATION is not defined.").Error())
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configFileLocation)

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.New(errors.New("Fatal error reading config file: ").Error() + err.Error())
	}

	for _, telegram_bot_variable := range telegramBotVariables {
		if !viper.IsSet("telegram_bot." + telegram_bot_variable) {
			return config, errors.New("Fatal error config: no telegram_bot " + telegram_bot_variable + " was found.")
		}
	}

	telegrambotConfig := TelegramBot{Token: viper.GetString("telegram_bot.token")}

	config.TelegramBot = telegrambotConfig

	return config, nil
}
