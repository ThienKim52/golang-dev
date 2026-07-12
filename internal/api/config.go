package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

// Config struct for app config
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"health-check-service"`
	InstanceID  string `envconfig:"INSTANCE_ID" default:""`
	Port        string `envconfig:"APP_PORT" default:"8080"`
	RedisAddr   string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisDB     int    `envconfig:"REDIS_DB" default:"0"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	// Generate UUID if INSTANCE_ID is empty
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
