package util

import (
	"bytes"
	"encoding/binary"
)

func Uint64FromVarByteArray(in []byte) uint64 {
	var corr []byte
	if len(in) < 8 {
		corr = append(bytes.Repeat([]byte{0x00}, 8-len(in)), in...)
		return binary.BigEndian.Uint64(corr)
	}
	return binary.BigEndian.Uint64(in)
}
