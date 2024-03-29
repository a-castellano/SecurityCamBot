package alarmsensors

import (
	"fmt"
	"strings"

	config "github.com/a-castellano/SecurityCamBot/config_reader"
	"github.com/streadway/amqp"
	tb "gopkg.in/tucnak/telebot.v2"
)

func ReceiveAlarmSensorMessages(rabbitmqConfig config.Rabbitmq, allowedSenders map[int]config.TelegramAllowedSender, messageQueue string, bot *tb.Bot) error {

	conn, errDial := amqp.Dial(rabbitmqConfig.GetDial())
	defer conn.Close()

	if errDial != nil {
		return errDial
	}

	channel, errChannel := conn.Channel()
	defer channel.Close()
	if errChannel != nil {
		return errChannel
	}

	_, errQueue := channel.QueueDeclare(
		messageQueue,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if errQueue != nil {
		return errQueue
	}

	errChannelQos := channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if errChannelQos != nil {
		return errChannelQos
	}

	messagesToProcess, errMessagesToProcess := channel.Consume(
		messageQueue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if errMessagesToProcess != nil {
		return errMessagesToProcess
	}

	processMessages := make(chan bool)

	go func() {
		for message := range messagesToProcess {

			notification := fmt.Sprintf("%s", message.Body)
			message.Ack(false)
			for _, userToNotify := range allowedSenders {
				sendMessage := false
				if strings.HasPrefix(notification, "DEBUG") {
					if userToNotify.SendDebug {
						sendMessage = true
					}
				} else {
					sendMessage = true
				}
				if sendMessage {
					user := &tb.User{ID: int64(userToNotify.ID)}
					bot.Send(user, notification)
				}
			}
		}
		return
	}()

	<-processMessages

	return nil
}
