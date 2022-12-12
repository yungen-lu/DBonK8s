package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type (
	// Config -.
	Config struct {
		App
		HTTP
		Log
		Line
		// PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME" env-default:"not set"`
		Version string `env:"APP_VERSION" env-default:"not set"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"3000"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"INFO"`
	}
	// Line -.
	Line struct {
		Token  string `env:"ACCESS_TOKEN" env-required:"true"`
		Secret string `env:"SECRET" env-required:"true"`
	}

	// PG -.
	// PG struct {
	// 	PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
	// 	URL     string `env-required:"true"                 env:"PG_URL"`
	// }
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// err := cleanenv.ReadConfig("./config/config.yml", cfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("config error: %w", err)
	// }
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Info("can't find .env will use env var")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
