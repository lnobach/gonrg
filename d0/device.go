//go:build !gonrgmocks

package d0

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	D0Timeout = 4 * time.Second
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

	if c.BaudRate <= 0 {
		c.BaudRate = 9600
	}

	return nil

}

func (d *deviceImpl) Get() (string, error) {

	sermode := &serial.Mode{
		BaudRate: d.config.BaudRate,
		DataBits: 7,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
		InitialStatusBits: &serial.ModemOutputBits{
			RTS: false,
			DTR: false,
		},
	}

	port, err := serial.Open(d.config.Device, sermode)
	if err != nil {
		return "", err
	}
	defer port.Close()
	err = port.SetReadTimeout(D0Timeout)
	if err != nil {
		return "", err
	}

	_, err = port.Write([]byte("/?!\x0D\x0A"))
	if err != nil {
		return "", err
	}
	err = port.Drain()
	if err != nil {
		return "", err
	}

	response := ""

	foundStart := false
	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		if !foundStart && strings.HasPrefix(line, "/") {
			foundStart = true
		}
		if foundStart {
			if strings.TrimSpace(line) == "!" {
				break
			}
			response += line + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		return response, err
	}

	return response, nil

}
