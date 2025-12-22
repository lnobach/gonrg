package d0

import (
	"time"

	"github.com/lnobach/gonrg/obis"
)

type Device interface {

	// Fetches raw data over the D0 interface
	Get() (ParseableRawData, error)
}

type ParseableRawData interface {
	ParseObis(cfg *ParseConfig, foundSet func(*obis.OBISEntry) error) (string, error)
}

type Parser interface {
	GetOBISMap(data ParseableRawData, measurementTime time.Time) (*obis.OBISMappedResult, error)
	GetOBISList(data ParseableRawData, measurementTime time.Time) (*obis.OBISListResult, error)
}
