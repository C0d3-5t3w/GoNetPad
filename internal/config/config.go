package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WebSocketPort string `yaml:"webSocketPort"`
	HTMLPort      string `yaml:"htmlPort"`
	WindowTitle   string `yaml:"windowTitle"`
	WindowWidth   int    `yaml:"windowWidth"`
	WindowHeight  int    `yaml:"windowHeight"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
}

func init() {
	configPath := "pkg/config/config.yaml"
	LoadConfig(configPath)
}
