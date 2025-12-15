package d0

import "github.com/lnobach/gonrg/options"

type DeviceConfig struct {

	// Device to connect to
	Device string

	// If 0, choose the recommended baud rate for D0 and your options.
	BaudRate int

	// Options chosen for the device, e.g. for compatibility with other meters.
	DeviceOptions options.Options
}

type ParseConfig struct {
	StrictMode bool

	// Options chosen for the device, e.g. for compatibility with other meters.
	DeviceOptions options.Options
}
