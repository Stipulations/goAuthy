package config

import(
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppState   string `env:"AppState,required"`
	ListenAddr string `env:"ListenerAddr,required"`
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}