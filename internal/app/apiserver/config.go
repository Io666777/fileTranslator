package apiserver

type Config struct {
	BindAddr          string `toml:"bind_addr"`
	LogLevel          string `toml:"log_level"`
	DatabaseURL       string `toml:"database_url"`
	SessionKey        string `toml:"session_key"`
	LibreTranslateURL string `toml:"libretranslate_url"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":5500",
		LogLevel: "debug",
	}
}
