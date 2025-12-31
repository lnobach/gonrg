//go:build !gonrgmocks

package sml

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/util"
	log "github.com/sirupsen/logrus"
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

	if c.ReconnectPause <= 0 {
		c.ReconnectPause = 5 * time.Second
	}

	return nil

}

func (d *deviceImpl) getSerMode() *serial.Mode {
	return &serial.Mode{
		BaudRate: d.config.BaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		InitialStatusBits: &serial.ModemOutputBits{
			RTS: false,
			DTR: false,
		},
	}
}

func (d *deviceImpl) setupSerialPort() (serial.Port, error) {

	sermode := d.getSerMode()

	port, err := serial.Open(d.config.Device, sermode)
	if err != nil {
		return nil, err
	}
	err = port.SetReadTimeout(8 * time.Second)
	if err != nil {
		return nil, err
	}

	return port, nil
}

func (d *deviceImpl) Get() (d0.ParseableRawData, error) {

	port, err := d.setupSerialPort()
	if err != nil {
		return nil, err
	}
	defer util.LogDeferWarn(port.Close)

	msg, err := GetNextRaw(port)
	if err != nil {
		return nil, err
	}

	return RawDataFromBytes(msg), nil

}

func (d *deviceImpl) GetForever(ctx context.Context, rcv chan d0.ParseableRawData) {

	chRaw := make(chan []byte)

	go func() {
		for {
			select {
			case rcv <- RawDataFromBytes(<-chRaw):
			case <-ctx.Done():
				log.WithError(ctx.Err()).Debugf("stopped parseable data pipe")
				return
			}
		}
	}()

	for {
		err := d.getForeverUntilErr(ctx, chRaw)
		if err == nil {
			break
		}
		log.WithError(err).Errorf("error while trying to get from serial port, retrying in %s",
			d.config.ReconnectPause)
		time.Sleep(d.config.ReconnectPause)
	}

}

func (d *deviceImpl) getForeverUntilErr(ctx context.Context, rcv chan []byte) (err error) {
	defer func() {
		err = util.PanicToError(recover(), err)
	}()
	port, err2 := d.setupSerialPort()
	if err2 != nil {
		err = fmt.Errorf("error while setting up serial port: %w", err2)
		return
	}
	defer util.LogDeferWarn(port.Close)
	err2 = GetForeverRaw(ctx, port, rcv)
	if err2 != nil {
		err = fmt.Errorf("error while getting from serial port: %w", err2)
	}
	return
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
