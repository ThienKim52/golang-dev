package redis

import (
	"github.com/kelseyhightower/envconfig"
)

// Config struct for app config
type config struct {
	Address  string `default:"Localhost:6379" envconfig:"address"`
	Password string `default:"" envconfig:"password"`
	DB       int    `default:"0" envconfig:"db"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
