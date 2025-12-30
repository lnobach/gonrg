package sml

import (
	"bufio"
	"bytes"
	"context"
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

func GetForeverRaw(ctx context.Context, r io.Reader, rcv chan []byte) error {

	bufr := bufio.NewReader(r)

	_, err := readUntil(bufr, StartSeq)
	if err != nil {
		return err
	}

	for {

		msg, err := readUntil(bufr, StartSeq)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Debugf("stopped raw data reading")
			return nil
		case rcv <- msg:
		default:
			log.Warn("read a new sml message but receiver can't keep up. Dropping message.")
		}

	}

}

func readUntil(r *bufio.Reader, seq []byte) ([]byte, error) {
	buf := make([]byte, 0, 256)
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
