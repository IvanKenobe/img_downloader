package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	GRPC GRPCConfig `yaml:"grpc" env-required:"true"`
	Nats NatsConfig `yaml:"nats" env-required:"true"`
}

type GRPCConfig struct {
	Host    string        `yaml:"host" env-required:"true"`
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type NatsConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()

	if path == "" {
		panic("config file path is empty")
	}

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		panic("config file not found" + path)
	}

	err = cleanenv.ReadConfig(path, &cfg)

	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")

	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
