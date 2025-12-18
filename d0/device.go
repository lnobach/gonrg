//go:build !gonrgmocks

package d0

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
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

	if c.D0Timeout <= 0 {
		c.D0Timeout = 8 * time.Second
	}

	log.Debugf("Using device option(s) %s", c.DeviceOptions)

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

	start := time.Now()

	port, err := serial.Open(d.config.Device, sermode)
	if err != nil {
		return "", err
	}
	defer port.Close()
	err = port.SetReadTimeout(d.config.D0Timeout)
	if err != nil {
		return "", err
	}

	if d.config.DeviceOptions.HasOption("0preamble") {

		log.Debugf("sending 0preamble...")

		_, err = port.Write([]byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
		if err != nil {
			return "", err
		}
		_, err = port.Write([]byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
		if err != nil {
			return "", err
		}

	}

	log.Debugf("sending request code...")
	_, err = port.Write([]byte("/?!\x0D\x0A"))
	if err != nil {
		return "", err
	}
	err = port.Drain()
	if err != nil {
		return "", err
	}

	if d.config.ResponseDelay > 0 {
		log.Debugf("waiting response_delay...")
		time.Sleep(d.config.ResponseDelay)
	}

	response := ""

	log.Debugf("reading response...")
	baudRateChanged := false
	foundStart := false
	scanner := bufio.NewScanner(NewTimeoutReader(port))
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
			if !baudRateChanged {
				if d.config.BaudRateRead > 0 {
					sermode.BaudRate = d.config.BaudRateRead
					port.SetMode(sermode)
				}

				baudRateChanged = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return response, err
	}
	log.Debugf("completed reading of response. Total time for d0 transaction: %s",
		time.Since(start))

	return response, nil

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
