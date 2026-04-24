package app

import "os"

type Config struct {
	Env string
}

func NewConfig() Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	return Config{Env: env}
}
