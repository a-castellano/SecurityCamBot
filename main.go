package main

import (
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

	bot, err := tb.NewBot(tb.Settings{
		Token:  botConfig.TelegramBot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/hello", func(m *tb.Message) {
		sender_id := m.Sender.ID
		log.Println("/hello received from sender", sender_id, ".")
		bot.Send(m.Sender, "Hello World!")
	})

	bot.Start()
}
