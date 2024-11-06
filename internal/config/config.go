package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-default:"50051"`
}

func New(path string) *Config {
	cfg := Config{}

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
