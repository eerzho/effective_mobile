package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env-required:"true"`
	Http     `yaml:"http"`
	Postgres `yaml:"postgres"`
	Api      `yaml:"api"`
}

type Http struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	Url            string `yaml:"url" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
}

type Api struct {
	ApiForAge string `yaml:"api_for_age"`
	ApiForGen string `yaml:"api_for_gen"`
	ApiForNat string `yaml:"api_for_nat"`
}

func MustLoad() *Config {
	// from env
	log.Print("Starting to get CONFIG_PATH variable")
	path, err := fetchPathEnv()
	if err == nil {
		return MustLoadByPath(path)
	}
	log.Print(err.Error())

	panic("config path is empty")
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &config
}

func fetchPathEnv() (string, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return "", fmt.Errorf("CONFIG_PATH variable is not set")
	}

	return path, nil
}
