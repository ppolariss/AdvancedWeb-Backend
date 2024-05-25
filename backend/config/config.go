package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
)

var Config struct {
	MossUrl    string `env:"MOSS_URL"`
	MossApiKey string `env:"MOSS_API_KEY"`
	RedisUrl   string `env:"REDIS_URL"`
}

func InitConfig() (err error) {
	if err = env.Parse(&Config); err != nil {
		return err
	}
	fmt.Printf("%+v\n", &Config)
	return nil
}
