package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func DecimalScaleToString(value int64, scale int) string {
	num := fmt.Sprintf("%d", value)

	if -scale >= len(num) {
		var isnegative bool
		if string(num[0]) == "-" {
			num = num[1:]
			isnegative = true
		}
		for -scale >= len(num) {
			num = "0" + num
		}
		if isnegative {
			num = "-" + num
		}
	}

	if scale < 0 {
		scaleRev := len(num) + scale
		num = num[:scaleRev] + "." + num[scaleRev:]
	} else {
		for i := 0; i < scale; i++ {
			num = num + "0"
		}
	}

	return num
}

func BytesToUint64(in []byte) uint64 {
	var corr []byte
	if len(in) < 8 {
		corr = append(bytes.Repeat([]byte{0x00}, 8-len(in)), in...)
		return binary.BigEndian.Uint64(corr)
	}
	return binary.BigEndian.Uint64(in)
}

func BytesToInt64(b []byte) int64 {

	var v int64

	for _, by := range b {
		v = (v << 8) | int64(by)
	}

	// Sign extension if the highest bit is set
	shift := uint(64 - len(b)*8)
	v = (v << shift) >> shift

	return v
}
