package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-default:"./storage/sso.db"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port        int           `yaml:"port" env-default:"44044"`
	Timeout     time.Duration  `yaml:"timeout" env-default:"30s"`
}

func MustLoad() *Config {
	 path := fetchConfigPath()
	 if path == "" {
		panic("config path is not specified")
	 }

	 if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	 }

	 var cfg Config
	 if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	 }

	 return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "Path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}