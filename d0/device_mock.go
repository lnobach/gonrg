//go:build gonrgmocks

package d0

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lnobach/gonrg/util"
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

func (d *deviceImpl) Get() (ParseableRawData, error) {
	response, exists := mock_map[d.config.Device]
	if !exists {
		return nil, fmt.Errorf("device not found: %s, available mock devices: %v", d.config.Device, util.KeysToStr(mock_map, ","))
	}
	log.Warnf("not using a real device, mocking %s", d.config.Device)

	if d.config.ResponseDelay > 0 {
		log.Debugf("waiting response_delay...")
		time.Sleep(d.config.ResponseDelay)
	}

	return RawDataFromString(response), nil
}

func (d *deviceImpl) GetForever(ctx context.Context, rcv chan ParseableRawData) {
	panic("GetForever unimplemented for plain D0")
}
