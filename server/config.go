package server

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/lnobach/gonrg/d0"
)

type ServerConfig struct {
	ListenAddr     string
	TrustedProxies []string

	// Device to connect to
	Meters []*ServedMeterConfig
}

type ServedMeterConfig struct {
	Name            string
	SML             bool
	Device          d0.DeviceConfig
	Parser          d0.ParseConfig
	RateLimitMaxAge time.Duration
	Cron            string
}

func ConfigFromFile(filename string) (*ServerConfig, error) {
	conf := &ServerConfig{}
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read server config file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, fmt.Errorf("could not read server config file: %w", err)
	}
	return conf, nil
}
