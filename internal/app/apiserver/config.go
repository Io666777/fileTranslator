package apiserver

import "github.com/Io666777/fileTranslator/internal/app/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":5500",
		LogLevel: "debug",
		Store: store.NewConfig(),
	}
}