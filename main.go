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

		senderID := upd.Message.Sender.ID
		if _, allowedSender := botConfig.TelegramBot.AllowedSenders[senderID]; !allowedSender {
			log.Println("Blocked message received from sender", senderID, ".")
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

	bot.Handle("/hello", func(m *tb.Message) {
		senderID := m.Sender.ID
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("/hello received from sender %s.", senderName)
		log.Println(logMsg)
		response := fmt.Sprintf("Hello %s.", senderName)
		bot.Send(m.Sender, response)
	})

	bot.Start()
}
