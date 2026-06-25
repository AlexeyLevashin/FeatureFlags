package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `env:"PORT" env-default:"8080"`

	DatabaseDSN string `env:"DATABASE_URL" env-required:"true"`
}

func LoadConfig() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err == nil {
		return &cfg
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Критическая ошибка загрузки конфига: %v", err)
	}

	return &cfg
}
