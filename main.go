package main

import (
	"fmt"
	"log"
	"log/syslog"
	"time"

	config "github.com/a-castellano/SecurityCamBot/config_reader"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

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
	startBotReplyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{rebootAllCamsBtn},
	}

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

		rebootCamReplyButtons := []tb.ReplyButton{}
		for webCamName, _ := range botConfig.Webcams {
			commandName := fmt.Sprintf("Reboot %s", webCamName)
			rebooCamBtn := tb.ReplyButton{Text: commandName}
			rebootCamReplyButtons = append(rebootCamReplyButtons, rebooCamBtn)
		}

		rebootCamReplyKeys := [][]tb.ReplyButton{rebootCamReplyButtons}
		response := "Select a cemera to be rebooted."
		bot.Send(m.Sender, response, &tb.ReplyMarkup{
			ReplyKeyboard: rebootCamReplyKeys,
		})
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Received text  unhandled message from sender %s.", senderName)
		fmt.Println(m.Text)
		log.Println(logMsg)
		response := fmt.Sprintf("Sorry %s, I don't know what are you talking about.", senderName)
		bot.Send(m.Sender, response)
	})

	bot.Start()
}
