package util_test

import (
	"testing"

	"github.com/lnobach/gonrg/util"
	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {

	assert.Equal(t, "12345", util.DecimalScaleToString(12345, 0))
	assert.Equal(t, "-12345", util.DecimalScaleToString(-12345, 0))
	assert.Equal(t, "1234.5678", util.DecimalScaleToString(12345678, -4))
	assert.Equal(t, "-1234.5678", util.DecimalScaleToString(-12345678, -4))
	assert.Equal(t, "0.00000012", util.DecimalScaleToString(12, -8))
	assert.Equal(t, "-0.00000012", util.DecimalScaleToString(-12, -8))
	assert.Equal(t, "10.000", util.DecimalScaleToString(10000, -3))
	assert.Equal(t, "-10.000", util.DecimalScaleToString(-10000, -3))

}

func TestByteToInt(t *testing.T) {

	assert.Equal(t, int64(-78851), util.BytesToInt64([]byte{0xfe, 0xcb, 0xfd}))

}
