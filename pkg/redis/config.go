package redis

import (
	"github.com/kelseyhightower/envconfig"
)

// Config struct for app config
type config struct {
	Address  string `default:"Localhost:6379" envconfig:"ADDRESS"`
	Password string `default:"" envconfig:"PASSWORD"`
	DB       int    `default:"0" envconfig:"DB"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
