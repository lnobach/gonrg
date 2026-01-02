package d0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	mt := time.Unix(12345678, 0)

	result, err := ParseOBISList(&ParseConfig{}, RawDataFromString(mock_ebz1str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 5)

	result, err = ParseOBISList(&ParseConfig{}, RawDataFromString(mock_ebz2str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 9)

	result, err = ParseOBISList(&ParseConfig{}, RawDataFromString(mock_ebz3str), mt)
	assert.NoError(t, err)
	assert.Equal(t, "EBZ5DD12345ETA_104", result.DeviceID)
	assert.Len(t, result.List, 9)

	result, err = ParseOBISList(&ParseConfig{}, RawDataFromString(mock_lugcuh50_1), mt)
	assert.NoError(t, err)
	assert.Equal(t, "LUGCUH50", result.DeviceID)
	assert.Len(t, result.List, 65)

	result, err = ParseOBISList(&ParseConfig{}, RawDataFromString(mock_lugcuh50_2), mt)
	assert.NoError(t, err)
	assert.Equal(t, "LUGCUH50", result.DeviceID)
	assert.Len(t, result.List, 65)

}
