package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"regexp"
	"strings"
	"time"

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

	rebootAllCamsBtn := tb.ReplyButton{Text: "📷  Reboot Cameras"}
	takeSnapshotBtn := tb.ReplyButton{Text: "📷  Take Snapshot"}
	startBotReplyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{rebootAllCamsBtn},
		[]tb.ReplyButton{takeSnapshotBtn},
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

	takeSnapshotFromCamReplyButtons := []tb.ReplyButton{}
	for webCamName := range botConfig.Webcams {
		commandName := fmt.Sprintf("From %s", webCamName)
		takeSnapshotFromCamBtn := tb.ReplyButton{Text: commandName}
		takeSnapshotFromCamReplyButtons = append(takeSnapshotFromCamReplyButtons, takeSnapshotFromCamBtn)
	}

	takeSnapshotFromCamReplyKeys := [][]tb.ReplyButton{takeSnapshotFromCamReplyButtons}

	//	takeSnapshotFromCamRegex := regexp.MustCompile(`From (.*)$`)

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
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Received text  unhandled message from sender %s.", senderName)
		log.Println(logMsg)
		stringText := string(m.Text)
		if strings.HasPrefix(stringText, "Reboot ") {

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

		} else {

			if strings.HasPrefix(stringText, "From ") {
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
						snapshotJobSuccessResponse := fmt.Sprintf("Snapshot job for webcam called '%s', has been sended.", camName)
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

			} else {
				response := fmt.Sprintf("Sorry %s, I don't know what are you talking about.", senderName)
				bot.Send(m.Sender, response,
					&tb.ReplyMarkup{
						ReplyKeyboard: startBotReplyKeys,
					})
			}
		}
	})

	go queues.ReceiveSnapshotJobs(botConfig.Rabbitmq, botConfig.Queues["receive_sanpshot"].Name, bot)
	bot.Start()
}
