package queues

import (
	"crypto/rand"
	"encoding/base32"

	config "github.com/a-castellano/SecurityCamBot/config_reader"
	jobs "github.com/a-castellano/WebCamSnapshotWorker/jobs"
	webcam "github.com/a-castellano/reolink-manager/webcam"
	"github.com/streadway/amqp"
	tb "gopkg.in/tucnak/telebot.v2"
)

func generateRandomString() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:32], nil
}

func SendJob(rabbitmqConfig config.Rabbitmq, jobQueue string, webcamInfo webcam.Webcam, senderID int) error {

	newJob := jobs.SnapshotJob{Errored: false, Finished: false, IP: webcamInfo.IP, User: webcamInfo.User, Password: webcamInfo.Password, Sender: senderID}

	newJob.ID, _ = generateRandomString()
	// For the time being Port and StreamPath are hardcoded
	newJob.Port = 554
	newJob.StreamPath = "/h264Preview_01_main"

	encodedJob, _ := jobs.EncodeJob(newJob)

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

	queue, errQueue := channel.QueueDeclare(
		jobQueue, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	if errQueue != nil {
		return errQueue
	}

	// send Job

	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	if err != nil {
		return err
	}
	return nil
}

func ReceiveSnapshotJobs(rabbitmqConfig config.Rabbitmq, jobQueue string, bot *tb.Bot) error {

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
		jobQueue,
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

	jobsToProcess, errJobsToProcess := channel.Consume(
		jobQueue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if errJobsToProcess != nil {
		return errJobsToProcess
	}

	processJobs := make(chan bool)

	go func() {
		for job := range jobsToProcess {

			decodedJob, decodeErr := jobs.DecodeJob(job.Body)
			if decodeErr == nil {
				snapshot := &tb.Photo{File: tb.FromDisk(decodedJob.SnapshotPath)}
				user := &tb.User{ID: int64(decodedJob.Sender)}
				bot.Send(user, snapshot)
				job.Ack(false)
			}
		}
		return
	}()

	<-processJobs

	return nil
}
