package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type ConnectConfig struct {
	HTTP_HOST string
	HTTP_PORT string
	DSN       string
}

type AppConfig struct {
	Read             time.Duration `yaml:"read"`
	Write            time.Duration `yaml:"write"`
	Idle             time.Duration `yaml:"idle"`
	ReadHeader       time.Duration `yaml:"read_header"`
	CheckHealth      time.Duration `yaml:"time_to_check_health"`
	DefaultTokenSize int64         `yaml:"default_token_size"`
}

type BackendConfig struct {
	BACKEND_SERVERS string `json:"servers"`
}

type Config struct {
	ConnectConfig
	AppConfig
	BackendConfig
}

func NewConfig(jsonPath, yamlPath string) (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Не удалось загрузить файл .env: %v", err)

		return Config{}, err
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Printf("Не удалось прочитать YAML-файл: %v", err)

		return Config{}, err
	}

	var appConfig AppConfig
	if err := yaml.Unmarshal(yamlFile, &appConfig); err != nil {
		log.Printf("Не удалось сделать Unmarshal в структуру: %v", err)

		return Config{}, err
	}

	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Printf("Не удалось прочитать json-файл: %v", err)

		return Config{}, err
	}

	var backendConfig BackendConfig
	if err := json.Unmarshal(jsonFile, &backendConfig); err != nil {
		log.Printf("Не удалось сделать Unmarshal в структуру: %v", err)

		return Config{}, err
	}

	return Config{
		ConnectConfig: ConnectConfig{
			HTTP_HOST: getEnv("HTTP_HOST"),
			HTTP_PORT: getEnv("HTTP_PORT"),
			DSN:       getEnv("PG_DSN"),
		},
		AppConfig:     appConfig,
		BackendConfig: backendConfig,
	}, nil
}
