package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var cfg *Config

// Can extend this to include DB config etc
type Config struct {
	PortNumber       string        `envconfig:"PORT_NUMBER"`
	ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT"`
	IsRunningInCloud bool          `envconfig:"IS_RUNNING_IN_CLOUD"`
}

func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		PortNumber:      ":21000",
		ShutdownTimeout: 5 * time.Second,
	}

	return cfg, envconfig.Process("", cfg)
}
