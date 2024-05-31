package config

import "github.com/spf13/viper"

//var Config struct {
//	Url string
//}

const (
	EnvUrl = "URL"
)

var defaultConfig = map[string]string{
	EnvUrl: "https://advanced-web-backend-service",
	// EnvUrl: "http://127.0.0.1:8080",
}

func init() {
	viper.AutomaticEnv()
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}
