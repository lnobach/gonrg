//go:build gonrgmocks

package sml

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
	log "github.com/sirupsen/logrus"
)

type deviceImpl struct {
	config *d0.DeviceConfig
	getCtr int64
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

	result := &MockMutableRawData{raw: msg, incr: d.getCtr}
	d.getCtr++
	return result, nil

}

type MockMutableRawData struct {
	raw  d0.ParseableRawData
	incr int64
}

func (d *deviceImpl) GetForever(ctx context.Context, rcv chan d0.ParseableRawData) {

	origMsg, err := RawDataFromFile(d.config.Device)
	if err != nil {
		log.WithError(err).Error("could not mock sml data from file: %w", err)
		return
	}

	log.Warnf("not using a real device, mocking from file %s", d.config.Device)

	for {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Debugf("stopped raw data reading")
			return
		case rcv <- &MockMutableRawData{raw: origMsg, incr: d.getCtr}:
		default:
			log.Warn("mocked a new sml message but receiver can't keep up. Dropping message.")
		}
		d.getCtr++

		time.Sleep(1 * time.Second)
	}

}

func (m *MockMutableRawData) ParseObis(cfg *d0.ParseConfig,
	foundSet func(*obis.OBISEntry) error) (string, error) {

	return m.raw.ParseObis(cfg, func(e *obis.OBISEntry) error {
		if e.ValueNum != 0 {
			if e.Unit == "kWh" {
				e.ValueNum += m.incr
			}
			if e.Unit == "W" {
				e.ValueNum += m.incr % 30
			}
		}
		return foundSet(e)
	})

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
