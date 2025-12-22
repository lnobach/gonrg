package sml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/lnobach/gonrg/util"
	log "github.com/sirupsen/logrus"
)

type TLVType uint8

type TLV struct {
	Type   uint8
	Value  []byte
	Elems  []*TLV
	Length int
	Depth  int
}

const (
	TLVType_OctetStream = 0
	TLVType_List        = 7
	TLVType_Unsigned    = 6
)

func TLVsFromBytes(buf []byte) ([]*TLV, error) {
	r := bytes.NewBuffer(buf)
	bufr := bufio.NewReader(r)
	return TLVsFromBuf(bufr)
}

func TLVsFromBuf(bufr *bufio.Reader) ([]*TLV, error) {
	var tlvs []*TLV
	for {
		tlv, err := tlvFromBuf(bufr, 0)
		if tlv != nil {
			tlvs = append(tlvs, tlv)
		}
		if err != nil || tlv == nil {
			remainder, rerr := io.ReadAll(bufr)
			if rerr != nil {
				log.WithError(rerr).Warnf("could not read remainder of sml packet")
				return tlvs, err
			}
			log.Debugf("remainder of msg is %x", remainder)
			return tlvs, err
		}
	}
}

func tlvFromBuf(r *bufio.Reader, depth int) (*TLV, error) {

	tlbyte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	t := tlbyte >> 4
	l := int(tlbyte & 0x0f)

	if depth == 0 && t == 0 && l == 0 {
		// end of messages
		return nil, nil
	}

	tlv := &TLV{
		Type:   t,
		Depth:  depth,
		Length: l,
	}

	if t == 7 { // list
		for i := 0; i < l; i++ {
			subtlv, err := tlvFromBuf(r, depth+1)
			if subtlv != nil {
				tlv.Elems = append(tlv.Elems, subtlv)
			}
			if err != nil {
				return tlv, err
			}
		}
		return tlv, nil
	}

	if l <= 0 {
		l = 1
	}

	tlv.Value = make([]byte, l-1)

	_, err = io.ReadFull(r, tlv.Value)
	if err != nil {
		return tlv, err
	}

	return tlv, nil

}

func (t *TLV) String() string {
	var desc string
	if t.Type == 7 {
		desc = fmt.Sprintf("Type: %d, Length: %d", t.Type, t.Length)
		for _, elem := range t.Elems {
			desc = fmt.Sprintf("%s\n%s", desc, elem)
		}
	} else {
		desc = fmt.Sprintf("Type: %d, Length: %d, Content %x, String %s", t.Type,
			t.Length, t.Value, util.BytesToPrintableString(t.Value))
	}
	for i := 0; i < t.Depth; i++ {
		desc = "  " + desc
	}
	return desc
}
