package config

import "github.com/spf13/viper"

//var Config struct {
//	Url string
//}

const (
	EnvUrl = "URL"
)

var defaultConfig = map[string]string{
	EnvUrl: "http://advanced-web-backend-service:8080",
	// EnvUrl: "http://127.0.0.1:8080",
}

func init() {
	viper.AutomaticEnv()
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}
