package d0

import (
	"time"

	"github.com/lnobach/gonrg/obis"
)

type Device interface {

	// Fetches raw data over the D0 interface
	Get() (string, error)
}

type Parser interface {
	GetOBISMap(rawdata string, measurementTime time.Time) (*obis.OBISMappedResult, error)
	GetOBISList(rawdata string, measurementTime time.Time) (*obis.OBISListResult, error)
}
