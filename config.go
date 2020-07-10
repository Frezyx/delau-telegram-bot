package main

import "flag"

// Config - структура описывающая конфиг сервера
type Config struct {
	BotToken string `toml:"bot_token"`
	AuthURL  string `toml:"auth_api_url"`
}

// NewConfig - создаем новый конфиг
func NewConfig() *Config {
	return &Config{
		BotToken: "null",
		AuthURL:  "null",
	}
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "botconfig.toml", "path to config file")
}
