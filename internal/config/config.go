package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig
	Storage DatabaseConfig
}

type ServerConfig struct {
	Endpoint string `yaml:"endpoint"`
}

type DatabaseConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	DBName   string        `yaml:"name"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config path is not set")
	}
	var config Config
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("config not read: %v", err)
	}
	return config
}
