package config

import (
	"errors"
	"reflect"

	webcam "github.com/a-castellano/Reolink-Manager/webcam"
	viperLib "github.com/spf13/viper"
)

type TelegramBot struct {
	Token          string
	AllowedSenders map[int]TelegramAllowedSender
}

type TelegramAllowedSender struct {
	Name string
	ID   int
}

type Config struct {
	TelegramBot TelegramBot
}

func contains(keys []string, keyName string) bool {
	for _, v := range keys {
		if v == keyName {
			return true
		}
	}

	return false
}

func ReadConfig() (Config, error) {
	var configFileLocation string
	var config Config

	var envVariable string = "SECURITY_CAM_BOT_CONFIG_FILE_LOCATION"

	requiredVariables := []string{"telegram_bot", "webcams"}
	telegramBotVariables := []string{"token", "allowed_senders"}
	allowedSendersVariables := []string{"name", "id"}
	webcamVariables := []string{"ip", "user", "password"}

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

	for _, requiredVariable := range requiredVariables {
		if !viper.IsSet(requiredVariable) {
			return config, errors.New("Fatal error config: no telegram_bot " + requiredVariable + " was found.")
		}
	}

	for _, telegramBotVariable := range telegramBotVariables {
		if !viper.IsSet("telegram_bot." + telegramBotVariable) {
			return config, errors.New("Fatal error config: no telegram_bot " + telegramBotVariable + " was found.")
		}
	}

	if !viper.IsSet("webcams.webcam") {
		return config, errors.New("Fatal error config: no webcams were found.")
	}

	senders := make(map[int]TelegramAllowedSender)

	readedNames := make(map[string]bool)
	readedIDs := make(map[int]bool)

	readedAllowedSenders := viper.GetStringMap("telegram_bot.allowed_senders")
	for sender_name, sender_info := range readedAllowedSenders {
		sender_info_value := reflect.ValueOf(sender_info)
		var newSender TelegramAllowedSender

		if sender_info_value.Kind() != reflect.Map {
			return config, errors.New("Fatal error config: allowed sender " + sender_name + " not a map.")
		} else {
			sender_info_value_map := sender_info_value.Interface().(map[string]interface{})
			keys := make([]string, 0, len(sender_info_value_map))
			for key_name := range sender_info_value_map {
				keys = append(keys, key_name)
			}
			for _, required_sender_key := range allowedSendersVariables {
				if !contains(keys, required_sender_key) {
					return config, errors.New("Fatal error config: allowed sender " + sender_name + " has no " + required_sender_key + ".")
				} else {
					if required_sender_key == "id" {
						newSender.ID = int(reflect.ValueOf(sender_info_value_map[required_sender_key]).Interface().(int64))
						if _, ok := readedIDs[newSender.ID]; ok {
							return config, errors.New("Fatal error config: allowed sender " + sender_name + " id is repeated.")
						} else {
							readedIDs[newSender.ID] = true
						}
					} else {
						if required_sender_key == "name" {
							newSender.Name = reflect.ValueOf(sender_info_value_map[required_sender_key]).Interface().(string)
							if _, ok := readedNames[newSender.Name]; ok {
								return config, errors.New("Fatal error config: allowed sender " + sender_name + " name is repeated.")
							} else {
								readedNames[newSender.Name] = true
							}
						}
					}
				}
			}
			senders[newSender.ID] = newSender
		}
	}

	webcams := make(map[string]webcam.Webcam)
	readedWebCamNames := make(map[string]bool)
	readedWebCamIPs := make(map[string]bool)
	readedWebcams := viper.GetStringMap("telegram_bot.webcams")
	for webcamName, webcamInfo := range readedWebcams {
		webcamInfoValue := reflect.ValueOf(webcamInfo)
		var newWebcam webcam.Webcam
		if webcamInfoValue.Kind() != reflect.Map {
			return config, errors.New("Fatal error config: webcam " + webcamName + " not a map.")
		} else {

			webcamInfoValueMap := webcamInfoValue.Interface().(map[string]interface{})

			keys := make([]string, 0, len(senderInfoValueMap))
			for key_name := range senderInfoValueMap {
				keys = append(keys, key_name)
			}
		}

	}

	telegrambotConfig := TelegramBot{Token: viper.GetString("telegram_bot.token"), AllowedSenders: senders}

	config.TelegramBot = telegrambotConfig

	return config, nil
}
