package config

import (
	"errors"
	"fmt"
	"net"
	"reflect"

	webcam "github.com/a-castellano/reolink-manager/webcam"
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

type Rabbitmq struct {
	Host     string
	Port     int
	User     string
	Password string
}

func (r Rabbitmq) GetDial() string {

	dialString := fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
	return dialString
}

type Queue struct {
	Name string
}

type Config struct {
	TelegramBot TelegramBot
	Webcams     map[string]webcam.Webcam
	Rabbitmq    Rabbitmq
	Queues      map[string]Queue
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

	requiredVariables := []string{"telegram_bot", "webcams", "rabbitmq"}
	telegramBotVariables := []string{"token", "allowed_senders"}
	allowedSendersVariables := []string{"name", "id"}
	webcamRequiredVariables := []string{"ip", "user", "password", "name"}
	rabbitmqRequiredVariables := []string{"host", "port", "user", "password"}
	rabbitmqRequiredQueues := []string{"send_sanpshot_commands", "receive_sanpshot"}

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
			return config, errors.New("Fatal error config: no " + requiredVariable + " field was found.")
		}
	}

	for _, telegramBotVariable := range telegramBotVariables {
		if !viper.IsSet("telegram_bot." + telegramBotVariable) {
			return config, errors.New("Fatal error config: no telegram_bot " + telegramBotVariable + " was found.")
		}
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
	readedWebCamIDs := make(map[string]bool)
	readedWebCamNames := make(map[string]bool)
	readedWebCamIPs := make(map[string]bool)
	readedWebcams := viper.GetStringMap("webcams")
	for webcamID, webcamInfo := range readedWebcams {
		webCamName := "NoName"
		webcamInfoValue := reflect.ValueOf(webcamInfo)
		var newWebcam webcam.Webcam
		if webcamInfoValue.Kind() != reflect.Map {
			return config, errors.New("Fatal error config: webcam " + webcamID + " not a map.")
		} else {

			if _, ok := readedWebCamIDs[webcamID]; ok {
				return config, errors.New("Fatal error config: webcam " + webcamID + " is repeated.")
			} else {

				webcamInfoValueMap := webcamInfoValue.Interface().(map[string]interface{})

				keys := make([]string, 0, len(webcamInfoValueMap))
				for key_name := range webcamInfoValueMap {
					keys = append(keys, key_name)
				}
				for _, requiredWebcamKey := range webcamRequiredVariables {
					if !contains(keys, requiredWebcamKey) {
						return config, errors.New("Fatal error config: webcam " + webcamID + " has no " + requiredWebcamKey + ".")
					} else {
						if requiredWebcamKey == "ip" {
							newWebcam.IP = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
							if net.ParseIP(newWebcam.IP) == nil {
								return config, errors.New("Fatal error config: webcam " + webcamID + " ip is invalid.")
							} else {
								if _, ok := readedWebCamIPs[newWebcam.IP]; ok {
									return config, errors.New("Fatal error config: webcam " + webcamID + " ip is repeated.")
								} else {
									readedWebCamIPs[newWebcam.IP] = true
								}
							}
						} else {
							if requiredWebcamKey == "user" {
								newWebcam.User = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
							} else {
								if requiredWebcamKey == "password" {
									newWebcam.Password = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
								} else {
									if requiredWebcamKey == "name" {
										webCamName = reflect.ValueOf(webcamInfoValueMap[requiredWebcamKey]).Interface().(string)
										if _, ok := readedWebCamNames[webCamName]; ok {
											return config, errors.New("Fatal error config: webcam " + webCamName + " name is repeated.")
										} else {
											readedWebCamNames[webCamName] = true
										}
									}
								}
							}
						}
					}
				}

				webcams[webCamName] = newWebcam
			}
		}
	}
	if len(webcams) == 0 {
		return config, errors.New("Fatal error config: no webcams were found.")
	}

	for _, rabbitmqVariable := range rabbitmqRequiredVariables {
		if !viper.IsSet("rabbitmq." + rabbitmqVariable) {
			return config, errors.New("Fatal error config: no rabbitmq " + rabbitmqVariable + " was found.")
		}
	}

	for _, rabbitmqRequiredQueue := range rabbitmqRequiredQueues {
		if !viper.IsSet("queues." + rabbitmqRequiredQueue) {
			return config, errors.New("Fatal error config: queue " + rabbitmqRequiredQueue + " was not found.")
		} else {
			if !viper.IsSet("queues." + rabbitmqRequiredQueue + ".name") {
				return config, errors.New("Fatal error config: queue " + rabbitmqRequiredQueue + " has no name.")
			}
		}
	}

	telegrambotConfig := TelegramBot{Token: viper.GetString("telegram_bot.token"), AllowedSenders: senders}

	rabbitmqConfig := Rabbitmq{Host: viper.GetString("rabbitmq.host"), Port: viper.GetInt("rabbitmq.port"), User: viper.GetString("rabbitmq.user"), Password: viper.GetString("rabbitmq.password")}

	config.TelegramBot = telegrambotConfig
	config.Webcams = webcams
	config.Rabbitmq = rabbitmqConfig
	queues := make(map[string]Queue)
	for _, rabbitmqRequiredQueue := range rabbitmqRequiredQueues {
		queue := Queue{Name: viper.GetString("queues." + rabbitmqRequiredQueue + ".name")}
		queues[rabbitmqRequiredQueue] = queue
	}
	config.Queues = queues

	return config, nil
}
