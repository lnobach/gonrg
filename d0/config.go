package d0

import (
	"time"

	"github.com/lnobach/gonrg/options"
)

type DeviceConfig struct {

	// Device to connect to
	Device string

	// If 0, choose the recommended baud rate for D0 and your options.
	BaudRate int

	// If 0, keep the baud rate at BaudRate
	BaudRateRead int

	// if not 0, wait before expecting response
	ResponseDelay time.Duration

	// Read timeout of the D0 serial connection
	D0Timeout time.Duration

	// Options chosen for the device, e.g. for compatibility with other meters.
	DeviceOptions options.Options

	// For continuous operations, time to wait until reconnect
	ReconnectPause time.Duration
}

type ParseConfig struct {
	StrictMode bool

	// Options chosen for the device, e.g. for compatibility with other meters.
	DeviceOptions options.Options
}
