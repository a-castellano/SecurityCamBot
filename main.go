package main

import (
	"bytes"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"regexp"
	"strings"
	"time"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	"github.com/a-castellano/SecurityCamBot/alarmmanager"
	config "github.com/a-castellano/SecurityCamBot/config_reader"
	queues "github.com/a-castellano/SecurityCamBot/queues"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "security-cam-bot")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove date prefix
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	botConfig, errConfig := config.ReadConfig()
	fmt.Println(botConfig.TelegramBot.AllowedSenders)
	if errConfig != nil {
		log.Fatal(errConfig)
		return
	}

	botPoller := &tb.LongPoller{Timeout: 15 * time.Second}

	senderWhiteList := tb.NewMiddlewarePoller(botPoller, func(upd *tb.Update) bool {
		if upd.Message == nil {
			return true
		}

		senderID := int(upd.Message.Sender.ID)
		if _, allowedSender := botConfig.TelegramBot.AllowedSenders[senderID]; !allowedSender {
			logError := fmt.Sprintf("Blocked message received from sender %d.", senderID)
			log.Println(logError)
			return false
		}

		return true
	})

	bot, err := tb.NewBot(tb.Settings{
		Token:  botConfig.TelegramBot.Token,
		Poller: senderWhiteList,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	watcher := apiwatcher.APIWatcher{Host: botConfig.AlarmManager.Host, Port: botConfig.AlarmManager.Port}
	alarmManagerRequester := apiwatcher.Requester{Client: client}

	apiInfo, apiInfoErr := watcher.ShowInfo(alarmManagerRequester)
	if apiInfoErr != nil {
		log.Fatal(apiInfoErr)
	}

	rebootAllCamsBtn := tb.ReplyButton{Text: "📷  Reboot Cameras"}
	takeSnapshotBtn := tb.ReplyButton{Text: "📷  Take Snapshot"}
	manageAlarmBtn := tb.ReplyButton{Text: "🔔  Alarm"}
	startBotReplyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{rebootAllCamsBtn},
		[]tb.ReplyButton{takeSnapshotBtn},
		[]tb.ReplyButton{manageAlarmBtn},
	}

	rebootCamReplyButtons := []tb.ReplyButton{}
	for webCamName := range botConfig.Webcams {
		commandName := fmt.Sprintf("Reboot %s", webCamName)
		rebooCamBtn := tb.ReplyButton{Text: commandName}
		rebootCamReplyButtons = append(rebootCamReplyButtons, rebooCamBtn)
	}

	rebootCamReplyKeys := [][]tb.ReplyButton{rebootCamReplyButtons}

	rebootCamRegex := regexp.MustCompile(`Reboot (.*)$`)
	snapshotCamRegex := regexp.MustCompile(`From (.*)$`)
	manageAlarmRegex := regexp.MustCompile(`Manage alarm (.*)$`)
	checkAlarmStatusRegex := regexp.MustCompile(`Check (.*) status\.$`)
	changeAlarmStatusRegex := regexp.MustCompile(`Change (.*) status\.$`)
	setAlarmStatusRegex := regexp.MustCompile(`Set (.*) to (.*)\.$`)

	takeSnapshotFromCamReplyButtons := []tb.ReplyButton{}
	for webCamName := range botConfig.Webcams {
		commandName := fmt.Sprintf("From %s", webCamName)
		takeSnapshotFromCamBtn := tb.ReplyButton{Text: commandName}
		takeSnapshotFromCamReplyButtons = append(takeSnapshotFromCamReplyButtons, takeSnapshotFromCamBtn)
	}

	takeSnapshotFromCamReplyKeys := [][]tb.ReplyButton{takeSnapshotFromCamReplyButtons}

	manageAlarmReplyButtons := []tb.ReplyButton{}
	alarmNameToIDMap := make(map[string]string)
	for alarmID, alarmInfo := range apiInfo.DevicesInfo {
		alarmNameToIDMap[alarmInfo.Name] = alarmID
		commandName := fmt.Sprintf("Manage alarm %s", alarmInfo.Name)
		manageAlarmBtn := tb.ReplyButton{Text: commandName}
		manageAlarmReplyButtons = append(manageAlarmReplyButtons, manageAlarmBtn)
	}

	manageAlarmReplyKeys := [][]tb.ReplyButton{manageAlarmReplyButtons}

	bot.Handle("/hello", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("/hello received from sender %s.", senderName)
		log.Println(logMsg)
		response := fmt.Sprintf("Hello %s.", senderName)
		bot.Send(m.Sender, response)
	})

	bot.Handle("/start", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("/start command received from sender %s.", senderName)
		log.Println(logMsg)
		response := fmt.Sprintf("Hello %s, please select and option.", senderName)

		bot.Send(m.Sender, response, &tb.ReplyMarkup{
			ReplyKeyboard: startBotReplyKeys,
		})
	})

	bot.Handle("📷  Reboot Cameras", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Manage Cameras command received from sender %s.", senderName)
		log.Println(logMsg)

		response := "Select a camera to be rebooted."
		bot.Send(m.Sender, response, &tb.ReplyMarkup{
			ReplyKeyboard: rebootCamReplyKeys,
		})
	})

	bot.Handle("🔔  Alarm", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Manage Alarm command received from sender %s.", senderName)
		log.Println(logMsg)

		response := "Select an Alarm to be manage."
		bot.Send(m.Sender, response, &tb.ReplyMarkup{
			ReplyKeyboard: manageAlarmReplyKeys,
		})
	})

	bot.Handle("📷  Take Snapshot", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Take snapshot from cameras command received from sender %s.", senderName)
		log.Println(logMsg)

		response := "From which camera?"
		bot.Send(m.Sender, response, &tb.ReplyMarkup{
			ReplyKeyboard: takeSnapshotFromCamReplyKeys,
		})
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		var textMatch bool = false
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Received text  unhandled message from sender %s.", senderName)
		log.Println(logMsg)
		stringText := string(m.Text)
		if strings.HasPrefix(stringText, "Reboot ") {
			textMatch = true

			camName := rebootCamRegex.FindStringSubmatch(stringText)[1]
			if _, ok := botConfig.Webcams[camName]; ok {

				response := fmt.Sprintf("Rebooting cam called '%s'.", camName)
				bot.Send(m.Sender, response)
				webcam := botConfig.Webcams[camName]
				connectErr := webcam.Connect(client)
				if connectErr != nil {
					log.Println(connectErr)
					connectErrResponse := fmt.Sprintf("Cannot connect with Webcam called '%s'.", camName)
					bot.Send(m.Sender, connectErrResponse,
						&tb.ReplyMarkup{
							ReplyKeyboard: startBotReplyKeys,
						})
				} else {
					rebootErr := webcam.Reboot(client)
					if rebootErr != nil {
						cantRebotedResponse := fmt.Sprintf("Cannot Reboot Webcam called '%s'.", camName)
						bot.Send(m.Sender, cantRebotedResponse,
							&tb.ReplyMarkup{
								ReplyKeyboard: startBotReplyKeys,
							})
						log.Println(rebootErr)
					} else {
						rebotedResponse := fmt.Sprintf("Webcam called '%s' has been rebooted.", camName)
						bot.Send(m.Sender, rebotedResponse,
							&tb.ReplyMarkup{
								ReplyKeyboard: startBotReplyKeys,
							})
					}
				}

			} else {
				response := fmt.Sprintf("Sorry %s, there is no cam called '%s'.", senderName, camName)
				bot.Send(m.Sender, response,
					&tb.ReplyMarkup{
						ReplyKeyboard: startBotReplyKeys,
					})
			}

		}

		if strings.HasPrefix(stringText, "From ") {
			textMatch = true
			camName := snapshotCamRegex.FindStringSubmatch(stringText)[1]
			if targetWebcam, ok := botConfig.Webcams[camName]; ok {

				aboutToSnapshotResponse := fmt.Sprintf("About to send a Snapshot job for Webcam called '%s'.", camName)
				log.Println(aboutToSnapshotResponse)
				sendJobErr := queues.SendJob(botConfig.Rabbitmq, botConfig.Queues["send_sanpshot_commands"].Name, targetWebcam, senderID)
				if sendJobErr != nil {
					log.Println(sendJobErr)
					snapshotJobErrResponse := fmt.Sprintf("Cannot send snapshot job for webcam called '%s', error was '%s'.", camName, sendJobErr.Error())
					bot.Send(m.Sender, snapshotJobErrResponse,
						&tb.ReplyMarkup{
							ReplyKeyboard: startBotReplyKeys,
						})

				} else {
					snapshotJobSuccessResponse := fmt.Sprintf("Snapshot job for webcam called '%s', has been sent.", camName)
					log.Println(snapshotJobSuccessResponse)
					bot.Send(m.Sender, snapshotJobSuccessResponse,
						&tb.ReplyMarkup{
							ReplyKeyboard: startBotReplyKeys,
						})

				}
			} else {
				response := fmt.Sprintf("Sorry %s, there is no cam called '%s'.", senderName, camName)
				bot.Send(m.Sender, response,
					&tb.ReplyMarkup{
						ReplyKeyboard: startBotReplyKeys,
					})
			}

		}

		if strings.HasPrefix(stringText, "Manage alarm ") {
			textMatch = true
			alarmName := manageAlarmRegex.FindStringSubmatch(stringText)[1]
			log.Println(alarmName, alarmNameToIDMap[alarmName])

			manageSelectedAlarmReplyButtons := []tb.ReplyButton{}
			manageSelectedAlarmBtnName := fmt.Sprintf("Check %s status.", alarmName)
			manageSelectedAlarmBtn := tb.ReplyButton{Text: manageSelectedAlarmBtnName}
			manageSelectedAlarmReplyButtons = append(manageSelectedAlarmReplyButtons, manageSelectedAlarmBtn)
			manageSelectedAlarmBtnChangeName := fmt.Sprintf("Change %s status.", alarmName)
			manageSelectedAlarmBtnChange := tb.ReplyButton{Text: manageSelectedAlarmBtnChangeName}
			manageSelectedAlarmReplyButtons = append(manageSelectedAlarmReplyButtons, manageSelectedAlarmBtnChange)
			response := "What do you want to do?"

			manageSelectedAlarmReplyKeys := [][]tb.ReplyButton{manageSelectedAlarmReplyButtons}
			bot.Send(m.Sender, response,
				&tb.ReplyMarkup{
					ReplyKeyboard: manageSelectedAlarmReplyKeys,
				})

		}

		if strings.HasPrefix(stringText, "Check ") {
			textMatch = true
			alarmName := checkAlarmStatusRegex.FindStringSubmatch(stringText)[1]
			log.Println(alarmName, alarmNameToIDMap[alarmName])
			apiInfo, apiInfoErr = watcher.ShowInfo(alarmManagerRequester)
			if apiInfoErr != nil {
				response := fmt.Sprintf("Error checking alarm status: %s", apiInfoErr.Error())
				bot.Send(m.Sender, response,
					&tb.ReplyMarkup{
						ReplyKeyboard: startBotReplyKeys,
					})

			} else {
				mode := apiInfo.DevicesInfo[alarmNameToIDMap[alarmName]].Mode
				if mode == "home" {
					mode = "Home Armed"
				}
				response := fmt.Sprintf("Alarm is %s.", mode)
				bot.Send(m.Sender, response,
					&tb.ReplyMarkup{
						ReplyKeyboard: startBotReplyKeys,
					})
			}

		}

		if strings.HasPrefix(stringText, "Change ") {
			textMatch = true
			alarmName := changeAlarmStatusRegex.FindStringSubmatch(stringText)[1]
			log.Println(alarmName, alarmNameToIDMap[alarmName])
			response := fmt.Sprintf("Set %s to:", alarmName)
			setSelectedAlarmReplyButtons := []tb.ReplyButton{}
			setSelectedAlarmDisarmedBtnName := fmt.Sprintf("Set %s to Disarmed.", alarmName)
			setSelectedAlarmDisarmedBtn := tb.ReplyButton{Text: setSelectedAlarmDisarmedBtnName}
			setSelectedAlarmReplyButtons = append(setSelectedAlarmReplyButtons, setSelectedAlarmDisarmedBtn)
			setSelectedAlarmArmedBtnName := fmt.Sprintf("Set %s to Armed.", alarmName)
			setSelectedAlarmArmedBtn := tb.ReplyButton{Text: setSelectedAlarmArmedBtnName}
			setSelectedAlarmReplyButtons = append(setSelectedAlarmReplyButtons, setSelectedAlarmArmedBtn)
			setSelectedAlarmHomeArmedBtnName := fmt.Sprintf("Set %s to HomeArmed.", alarmName)
			setSelectedAlarmHomeArmedBtn := tb.ReplyButton{Text: setSelectedAlarmHomeArmedBtnName}
			setSelectedAlarmReplyButtons = append(setSelectedAlarmReplyButtons, setSelectedAlarmHomeArmedBtn)

			setSelectedAlarmReplyKeys := [][]tb.ReplyButton{setSelectedAlarmReplyButtons}
			bot.Send(m.Sender, response,
				&tb.ReplyMarkup{
					ReplyKeyboard: setSelectedAlarmReplyKeys,
				})

		}

		if strings.HasPrefix(stringText, "Set ") {
			textMatch = true
			alarmName := setAlarmStatusRegex.FindStringSubmatch(stringText)[1]
			log.Println(alarmName, alarmNameToIDMap[alarmName])
			newMode := setAlarmStatusRegex.FindStringSubmatch(stringText)[2]
			log.Println(newMode)
			jsonString := fmt.Sprintf("{\"mode\":\"%s\"}", newMode)
			var jsonStr = []byte(jsonString)
			apiURL := fmt.Sprintf("http://%s:%d/devices/status/%s", botConfig.AlarmManager.Host, botConfig.AlarmManager.Port, alarmNameToIDMap[alarmName])
			req, _ := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			_, setNewModeErr := client.Do(req)

			var response string
			if setNewModeErr != nil {
				response = fmt.Sprintf("Error setting Alarm mode: %s", setNewModeErr.Error())
			} else {
				response = fmt.Sprintf("%s set to %s.", alarmName, newMode)
			}
			bot.Send(m.Sender, response,
				&tb.ReplyMarkup{
					ReplyKeyboard: startBotReplyKeys,
				})

		}

		if textMatch == false {
			response := fmt.Sprintf("Sorry %s, I don't know what are you talking about.", senderName)
			bot.Send(m.Sender, response,
				&tb.ReplyMarkup{
					ReplyKeyboard: startBotReplyKeys,
				})
		}
	})

	go queues.ReceiveSnapshotJobs(botConfig.Rabbitmq, botConfig.Queues["receive_sanpshot"].Name, bot)
	go alarmmanager.ReceiveAlarmMessages(botConfig.Rabbitmq, botConfig.TelegramBot.AllowedSenders, botConfig.Queues["alarmwatcher"].Name, bot)
	bot.Start()
}
