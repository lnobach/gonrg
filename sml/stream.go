package sml

import (
	"bufio"
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"
)

var (
	StartSeq = []byte{0x1b, 0x1b, 0x1b, 0x1b, 0x01, 0x01, 0x01, 0x01}
)

func GetNextRaw(r io.Reader) ([]byte, error) {

	bufr := bufio.NewReader(r)

	_, err := readUntil(bufr, StartSeq)
	if err != nil {
		return nil, err
	}

	msg, err := readUntil(bufr, StartSeq)
	if err != nil {
		return nil, err
	}

	log.Debugf("sml: obtained message, len=: %d, %x\n", len(msg), msg)

	return msg, nil

}

func readUntil(r *bufio.Reader, seq []byte) ([]byte, error) {
	var buf []byte
	for {
		s, err := r.ReadBytes(seq[len(seq)-1])
		if err != nil {
			return nil, err
		}

		buf = append(buf, s...)
		if bytes.HasSuffix(buf, seq) {
			return buf[:len(buf)-len(seq)], nil
		}
	}
}
