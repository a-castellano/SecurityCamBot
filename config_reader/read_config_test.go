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
	if config.TelegramBot.AllowedSenders[13].Name != "Bob" {
		t.Errorf("TelegramBot AllowedSenders with id 13 should be Bob. Returned: %s.", config.TelegramBot.AllowedSenders[12].Name)
	}
}
