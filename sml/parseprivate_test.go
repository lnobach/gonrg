package sml_test

import (
	"testing"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/sml"
	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	mt := time.Unix(12345678, 0)

	raw, err := sml.RawDataFromFile("../testdata/sml_priv/iskra.bin")
	assert.NoError(t, err)
	p, err := d0.NewParser(d0.ParseConfig{})
	assert.NoError(t, err)
	result, err := p.GetOBISMap(raw, mt)
	assert.NoError(t, err)
	assert.Equal(t, "..ISK....X", result.DeviceID)
	assert.Len(t, result.List, 5)
	assert.Equal(t, int64(3482638), result.List[2].ValueNum)
	assert.Equal(t, -4, result.List[2].ValueScale)

}
