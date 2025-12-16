//go:build gonrgmocks

package d0

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type deviceImpl struct {
	config *DeviceConfig
}

func NewDevice(config DeviceConfig) (Device, error) {
	c := &config
	err := deviceSetDefaults(c)
	if err != nil {
		return nil, fmt.Errorf("failure setting config: %w", err)
	}
	return &deviceImpl{config: c}, nil
}

func deviceSetDefaults(c *DeviceConfig) error {
	if strings.TrimSpace(c.Device) == "" {
		return fmt.Errorf("device must not be empty")
	}
	log.Debugf("Using device option(s) %s", c.DeviceOptions)
	return nil
}

func (d *deviceImpl) Get() (string, error) {
	response, exists := Mock_map[d.config.Device]
	if !exists {
		return "", fmt.Errorf("device not found: %s, available mock devices: %v", d.config.Device, keysToStr(Mock_map))
	}
	log.Warnf("not using a real device, mocking %s", d.config.Device)
	return response, nil
}

func keysToStr[T any](m map[string]T) string {
	if len(m) == 0 {
		return ""
	}
	str := ""
	for key, _ := range m {
		str += key + ","
	}
	return str[0 : len(str)-1]
}
