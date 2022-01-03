package queues

import (
	"crypto/rand"
	"encoding/base32"

	config "github.com/a-castellano/SecurityCamBot/config_reader"
	jobs "github.com/a-castellano/WebCamSnapshotWorker/jobs"
	webcam "github.com/a-castellano/reolink-manager/webcam"
	"github.com/streadway/amqp"
)

func generateRandomString() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:32], nil
}

func SendJob(rabbitmqConfig config.Rabbitmq, jobQueue string, webcamInfo webcam.Webcam) error {

	newJob := jobs.SnapshotJob{Errored: false, Finished: false, IP: webcamInfo.IP, User: webcamInfo.User, Password: webcamInfo.Password}

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
