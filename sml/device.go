//go:build !gonrgmocks

package sml

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/util"
	"go.bug.st/serial"
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

	sermode := &serial.Mode{
		BaudRate: d.config.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		InitialStatusBits: &serial.ModemOutputBits{
			RTS: false,
			DTR: false,
		},
	}

	port, err := serial.Open(d.config.Device, sermode)
	if err != nil {
		return nil, err
	}
	defer util.LogDeferWarn(port.Close)
	err = port.SetReadTimeout(8 * time.Second)
	if err != nil {
		return nil, err
	}

	msg, err := GetNextRaw(port)
	if err != nil {
		return nil, err
	}

	return RawDataFromBytes(msg), nil

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
