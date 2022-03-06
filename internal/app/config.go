package app

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

const configFile = "configs/config.yml"

// Config general configs for App
type Config struct {
	// URL to health check
	URL string `yaml:"url"`

	// retry after a some amount of time
	RetryAfter int64 `yaml:"retry_after"`

	Telegram telegram
}

// telegram
type telegram struct {
	Token  string
	ChatId string
}

// NewConfig is Config constructor
func NewConfig() (*Config, error) {
	arrByte, err := parseConfigFile()

	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = yaml.Unmarshal(arrByte, config)
	if err != nil {
		return nil, err
	}

	// load .env file
	err = loadEnvVariable()
	if err != nil {
		return nil, err
	}

	config.Telegram.Token = os.Getenv("telegram_token")
	config.Telegram.ChatId = os.Getenv("telegram_chat_id")

	configJSON, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		return nil, err
	}
	log.Println("Configuration:", string(configJSON))

	return config, nil
}

// parseConfigFile parse the config file
func parseConfigFile() ([]byte, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	arrByte, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return arrByte, nil
}

// loadEnvVariable load .env file
func loadEnvVariable() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}
