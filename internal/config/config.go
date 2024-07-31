package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func NewConfig() *Config {
	configPath := flag.String("config", "./config/local.yaml", "path to configuration yaml file")
	flag.Parse()

	// check if file exists
	_, err := os.Stat(*configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", *configPath)
	}

	var config Config

	fmt.Println(*configPath)

	if err := cleanenv.ReadConfig(*configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &config
}
