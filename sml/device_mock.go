//go:build gonrgmocks

package sml

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lnobach/gonrg/d0"
	log "github.com/sirupsen/logrus"
)

type deviceImpl struct {
	config *d0.DeviceConfig
}

func NewDevice(config d0.DeviceConfig) (*deviceImpl, error) {
	c := &config
	err := deviceSetDefaults(c)
	if err != nil {
		return nil, fmt.Errorf("failure setting config: %w", err)
	}
	return &deviceImpl{config: c}, nil
}

func deviceSetDefaults(c *d0.DeviceConfig) error {

	if strings.TrimSpace(c.Device) == "" {
		return fmt.Errorf("device must not be empty")
	}

	if c.BaudRate <= 0 {
		c.BaudRate = 9600
	}

	if c.D0Timeout <= 0 {
		c.D0Timeout = 8 * time.Second
	}

	return nil

}

func (d *deviceImpl) Get() (d0.ParseableRawData, error) {
	msg, err := RawDataFromFile(d.config.Device)
	if err != nil {
		return nil, fmt.Errorf("could not mock sml data from file: %w", err)
	}

	log.Warnf("not using a real device, mocking from file %s", d.config.Device)

	return msg, nil

}

type timeoutReader struct {
	r io.Reader
}

var errTimeout = errors.New("timed out while waiting for d0 response")

// see https://github.com/bugst/go-serial/issues/148
func NewTimeoutReader(r io.Reader) timeoutReader {
	return timeoutReader{r}
}

func (t timeoutReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n == 0 && err == nil {
		err = errTimeout
	}
	return
}
