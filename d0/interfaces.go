package d0

import (
	"context"

	"github.com/lnobach/gonrg/obis"
)

type Device interface {

	// Fetches raw data over the D0 interface
	Get() (ParseableRawData, error)

	GetForever(ctx context.Context, rcv chan ParseableRawData)
}

type ParseableRawData interface {
	ParseObis(cfg *ParseConfig, foundSet func(*obis.OBISEntry) error) (string, error)
}
