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

	var (
		// Universal markup builders.
		mainMenu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		//selector = &tb.ReplyMarkup{}

		// Reply buttons.
		btnCameraRebootOption = mainMenu.Text("📷  Reboot Cameras")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind,
		// as it has to be for callback routing to work.
		//
		//		btnPrev = selector.Data("⬅", "prev")
		//		btnNext = selector.Data("➡", "next")
	)

	mainMenu.Reply(
		mainMenu.Row(btnCameraRebootOption),
		//mainMenu.Row(btnSettings),
	)
	//	selector.Inline(
	//		selector.Row(btnPrev, btnNext),
	//	)

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
		bot.Send(m.Sender, response, mainMenu)
	})

	bot.Handle("📷  Reboot Cameras", func(m *tb.Message) {
		senderID := int(m.Sender.ID)
		senderName := botConfig.TelegramBot.AllowedSenders[senderID].Name
		logMsg := fmt.Sprintf("Manage Cameras command received from sender %s.", senderName)
		log.Println(logMsg)

		rebootWebcamsMenu := &tb.ReplyMarkup{}
		for webCamName, _ := range botConfig.Webcams {
			commandName := fmt.Sprintf("Reboot %s.", webCamName)
			rebootWebcamsMenu.Text(commandName)
		}

		response := "Select a cemera to be rebooted."
		mainMenu.ReplyKeyboardRemove()
		bot.Send(m.Sender, response, nil)
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
