package sml

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/util"

	log "github.com/sirupsen/logrus"
)

type RawData struct {
	Raw []byte
}

type obisscanner struct {
	cfg      *d0.ParseConfig
	foundSet func(*obis.OBISEntry) error
	deviceid string
}

var (
	_ d0.ParseableRawData = &RawData{}
)

func (p *RawData) ParseObis(cfg *d0.ParseConfig, foundSet func(*obis.OBISEntry) error) (string, error) {

	tlvs, err := TLVsFromBytes(p.Raw)
	if err != nil {
		return "", err
	}

	obisscanner := &obisscanner{
		cfg:      cfg,
		foundSet: foundSet,
	}

	err = obisscanner.scanForObisData(tlvs)
	if err != nil {
		return "", err
	}

	return obisscanner.deviceid, nil

}

func (s *obisscanner) scanForObisData(tlvs []*TLV) error {

	for _, root := range tlvs {
		if root.Type != TLVType_List {
			continue
		}
		if len(root.Elems) < 4 {
			continue
		}
		glr := root.Elems[3]
		if !isGetListResponse(glr) {
			continue
		}

		err := s.parseGetListResponseBody(glr.Elems[1])
		if err != nil {
			return fmt.Errorf("error parsing GetListResponse body: %w", err)
		}
		return nil

	}

	return errors.New("not found any obis data in sml tree")

}

func (s *obisscanner) parseGetListResponseBody(body *TLV) error {

	if len(body.Elems) < 5 {
		return errors.New("body does not have >= 5 elements")
	}

	serverid := body.Elems[1]
	if serverid.Type != TLVType_OctetStream {
		return errors.New("serverID is not of type octet-stream")
	}

	deviceID, err := s.parseDeviceID(serverid.Value)
	if err != nil {
		return fmt.Errorf("error parsing device ID: %w", err)
	}
	s.deviceid = deviceID

	return s.parseOBISList(body.Elems[4])

}

func (s *obisscanner) parseOBISList(list *TLV) error {
	if list.Type != TLVType_List {
		return errors.New("obis list tlv is not of list type")
	}

	for i, elem := range list.Elems {
		obisentry, err := ParseOBISTLV(elem, s.cfg)
		if err != nil {
			if s.cfg.StrictMode {
				return fmt.Errorf("could not parse obis element from tlv: %w", err)
			} else {
				log.WithError(err).Errorf("could not parse obis element at position %d, will continue", i)
				continue
			}
		}
		if obisentry != nil {
			err := s.foundSet(obisentry)
			if err != nil {
				if s.cfg.StrictMode {
					return fmt.Errorf("error addin obis entry: %w", err)
				} else {
					log.WithError(err).Errorf("could not add obis element at position %d, will continue", i)
					continue
				}
			}
		}
	}

	return nil
}

func (s *obisscanner) parseDeviceID(serverID []byte) (string, error) {

	// serverid parsing is completely based on analysis of the raw data
	// of different meters since no docs seem available.

	if len(serverID) < 1 {
		return "", nil
	}

	firstbyte := serverID[0]

	if firstbyte > 0x07 && firstbyte < 0x0f && len(serverID) >= 5 {
		prenum := serverID[1]
		vendorslug := util.BytesToPrintableString(serverID[2:5])
		serial := util.BytesToUint64(serverID[5:])
		return fmt.Sprintf("%d%s%010d", prenum, vendorslug, serial), nil
	}

	if firstbyte <= 0x07 && len(serverID) >= 4 {
		vendorslug := util.BytesToPrintableString(serverID[1:4])
		serial := util.BytesToUint64(serverID[5:])
		return fmt.Sprintf("%s%010d", vendorslug, serial), nil
	}

	if util.AllCharsPrintable(serverID) {
		return string(serverID), nil
	}

	return fmt.Sprintf("%x", serverID), nil

}

func isGetListResponse(tlv *TLV) bool {
	if tlv.Type != TLVType_List {
		return false
	}
	if len(tlv.Elems) < 2 {
		return false
	}
	tlv_type := tlv.Elems[0]
	if tlv_type.Type != TLVType_Unsigned {
		return false
	}
	return tlv_type.Value[len(tlv_type.Value)-2] == 0x07 &&
		tlv_type.Value[len(tlv_type.Value)-1] == 0x01
}

func RawDataFromBytes(raw []byte) *RawData {
	if idx := bytes.Index(raw, StartSeq); idx >= 0 {
		raw = raw[idx+len(StartSeq):]
		if idx := bytes.Index(raw, StartSeq); idx >= 0 {
			raw = raw[:idx]
		}
	}
	return &RawData{
		Raw: raw,
	}
}

func RawDataFromFile(file string) (*RawData, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return RawDataFromBytes(bytes), nil
}
