// +build integration_tests

package queues

import (
	"log"
	"testing"

	config "github.com/a-castellano/SecurityCamBot/config_reader"
	webcam "github.com/a-castellano/reolink-manager/webcam"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func TestSendJob(t *testing.T) {

	var rabbitConfig config.Rabbitmq

	rabbitConfig.Host = "rabbitmq"
	rabbitConfig.Port = 5672
	rabbitConfig.User = "guest"
	rabbitConfig.Password = "guest"

	queueName := "sendsnapshotjobs"

	targetWebcam := webcam.Webcam{IP: "10.10.10.10", User: "user", Password: "pass"}

	sendJobErr := SendJob(rabbitConfig, queueName, targetWebcam, 1)

	if sendJobErr != nil {
		t.Errorf("TestSendJob shouldn't fail, error was \"%s\"", sendJobErr.Error())
	}

}
