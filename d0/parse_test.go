package d0_test

import (
	"testing"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	mt := time.Unix(12345678, 0)

	result, err := d0.ParseOBISList(&d0.ParseConfig{}, d0.RawDataFromString(d0.Mock_ebz1str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 5)

	result, err = d0.ParseOBISList(&d0.ParseConfig{}, d0.RawDataFromString(d0.Mock_ebz2str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 9)

	result, err = d0.ParseOBISList(&d0.ParseConfig{}, d0.RawDataFromString(d0.Mock_ebz3str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 9)

	result, err = d0.ParseOBISList(&d0.ParseConfig{}, d0.RawDataFromString(d0.Mock_lugcuh50_1), mt)
	assert.NoError(t, err)
	assert.Equal(t, "LUGCUH50", result.DeviceID)
	assert.Len(t, result.List, 65)

	result, err = d0.ParseOBISList(&d0.ParseConfig{}, d0.RawDataFromString(d0.Mock_lugcuh50_2), mt)
	assert.NoError(t, err)
	assert.Equal(t, "LUGCUH50", result.DeviceID)
	assert.Len(t, result.List, 65)

}
