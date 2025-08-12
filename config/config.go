package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Auth  AuthConfig
	DB    DBConfig
	HTTP  HTTPConfig
	Redis RedisConfig
}

func NewConfig() *Config {
	err := godotenv.Load(".env", ".env_secrets")
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}

	return &Config{
		Auth:  LoadAuthConfig(),
		DB:    LoadDBConfig(),
		HTTP:  LoadHTTPConfig(),
		Redis: LoadRedisConfig(),
	}
}
