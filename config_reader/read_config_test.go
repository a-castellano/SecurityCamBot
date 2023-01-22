package config

import (
	"os"
	"strings"
	"testing"
)

func TestProcessNoConfigFilePresent(t *testing.T) {

	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without any valid config file should fail.")
	} else {
		if err.Error() != "Environment variable SECURITY_CAM_BOT_CONFIG_FILE_LOCATION is not defined." {
			t.Errorf("Error should be 'Environment variable SECURITY_CAM_BOT_CONFIG_FILE_LOCATION is not defined.', but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoTelegramToken(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_telegram_token/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without telegram token should fail.")
	} else {
		if err.Error() != "Fatal error config: no telegram_bot token was found." {
			t.Errorf("Error should be \"Fatal error config: no telegram_bot token was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoAllowedSenders(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_allowed_senders/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without allowed senders should fail.")
	} else {
		if err.Error() != "Fatal error config: no telegram_bot allowed_senders was found." {
			t.Errorf("Error should be \"Fatal error config: no telegram_bot allowed_senders was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigRepeatedAllowedSenders(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_senders/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with repeated allowed senders should fail.")
	} else {
		if !strings.Contains(err.Error(), "id is repeated.") {
			t.Errorf("Error should contain \"id is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigRepeatedAllowedSenders2(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_senders2/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with repeated allowed senders should fail.")
	} else {
		if !strings.Contains(err.Error(), "name is repeated.") {
			t.Errorf("Error should contain \"name is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCamsField(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cams_field/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcams info should fail.")
	} else {
		if err.Error() != "Fatal error config: no webcams field was found." {
			t.Errorf("Error should be \"Fatal error config: no webcams field was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCams(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cams/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcams info should fail.")
	} else {
		if err.Error() != "Fatal error config: no webcams were found." {
			t.Errorf("Error should be \"Fatal error config: no webcams were found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCamIP(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cam_ip/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcam IP should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam01 has no ip." {
			t.Errorf("Error should be \"Fatal error config: webcam cam01 has no ip.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCamUser(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cam_user/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcam user should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam01 has no user." {
			t.Errorf("Error should be \"Fatal error config: webcam cam01 has no user.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCamPassword(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cam_password/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcam password should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam01 has no password." {
			t.Errorf("Error should be \"Fatal error config: webcam cam01 has no password.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigInvalidWebCamIP(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_invalid_ip/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with invalid webcam IP should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam01 ip is invalid." {
			t.Errorf("Error should be \"Fatal error config: webcam cam01 ip is invalid.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWebCamRepetedIP(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_ips/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with repeated webcam IP should fail.")
	} else {
		if !strings.Contains(err.Error(), "ip is repeated.") {
			t.Errorf("Error should contain \"ip is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoWebCamName(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_web_cam_name/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with no webcam Name should fail.")
	} else {
		if err.Error() != "Fatal error config: webcam cam01 has no name." {
			t.Errorf("Error should be \"Fatal error config: webcam cam01 has no name.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigWebCamRepetedName(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_repeated_names/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method with repeated webcam name should fail.")
	} else {
		if !strings.Contains(err.Error(), "name is repeated.") {
			t.Errorf("Error should contain \"name is repeated.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoRabbitServer(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_rabbitmq/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without rabbitmq config should fail.")
	} else {
		if err.Error() != "Fatal error config: no rabbitmq field was found." {
			t.Errorf("Error should be \"Fatal error config: no rabbitmq field was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoRabbitUser(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_rabbitmq_user/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without rabbitmq user should fail.")
	} else {
		if err.Error() != "Fatal error config: no rabbitmq user was found." {
			t.Errorf("Error should be \"Fatal error config: no rabbitmq user was found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoQueueName(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_receive_sanpshots_queue_name/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without queue name should fail.")
	} else {
		if err.Error() != "Fatal error config: queue receive_sanpshot has no name." {
			t.Errorf("Error should be \"Fatal error config: queue receive_sanpshot has no name.\" but error was '%s'.", err.Error())
		}
	}
}

func TestProcessConfigNoQueue(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_no_receive_sanpshots_queue/")
	_, err := ReadConfig()
	if err == nil {
		t.Errorf("ReadConfig method without required queue should fail.")
	} else {
		if err.Error() != "Fatal error config: queue receive_sanpshot was not found." {
			t.Errorf("Error should be \"Fatal error config: queue receive_sanpshot was not found.\" but error was '%s'.", err.Error())
		}
	}
}

func TestOKConfig(t *testing.T) {
	os.Setenv("SECURITY_CAM_BOT_CONFIG_FILE_LOCATION", "./config_files_test/config_ok/")
	config, err := ReadConfig()
	if err != nil {
		t.Errorf("ReadConfig with ok config shouln't return errors. Returned: %s.", err.Error())
	}
	if config.TelegramBot.Token != "token" {
		t.Errorf("TelegramBot token should be token. Returned: %s.", config.TelegramBot.Token)
	}
	if len(config.TelegramBot.AllowedSenders) != 2 {
		t.Errorf("TelegramBot AllowedSenders length should be 2. Returned: %d.", len(config.TelegramBot.AllowedSenders))
	}
	if config.TelegramBot.AllowedSenders[12].Name != "Alice" {
		t.Errorf("TelegramBot AllowedSenders with id 12 should be Alice. Returned: %s.", config.TelegramBot.AllowedSenders[12].Name)
	}
	if config.TelegramBot.AllowedSenders[12].SendDebug != false {
		t.Errorf("TelegramBot AllowedSenders with id 12 SendDebug value should be false. Returned: true.")
	}
	if config.TelegramBot.AllowedSenders[13].Name != "Bob" {
		t.Errorf("TelegramBot AllowedSenders with id 13 should be Bob. Returned: %s.", config.TelegramBot.AllowedSenders[12].Name)
	}
	if config.TelegramBot.AllowedSenders[13].SendDebug != true {
		t.Errorf("TelegramBot AllowedSenders with id 13 SendDebug value should be true. Returned: false.")
	}
	if len(config.Webcams) != 2 {
		t.Errorf("Config should contain 2 webcams. Returned: %d.", len(config.Webcams))
	}
	if config.Webcams["cam2"].IP != "10.10.10.35" {
		t.Errorf("TelegramBot cam02 should have IP 10.10.10.35. Returned: %s.", config.Webcams["cam2"].IP)
	}
	if config.Rabbitmq.User != "guest" {
		t.Errorf("Rabbitmq user should be guest. Returned: %s.", config.Rabbitmq.User)
	}
	if config.Queues["send_sanpshot_commands"].Name != "incoming" {
		t.Errorf("Queue send_sanpshot_commands name should be incoming. Returned: %s.", config.Queues["send_sanpshot_commands"].Name)
	}
}
