package sml

import (
	"errors"
	"fmt"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
)

func ParseOBISTLV(list *TLV, parseConfig *d0.ParseConfig) (*obis.OBISEntry, error) {

	if list.Type != TLVType_List {
		return nil, errors.New("obis tlv is not of type list")
	}

	if len(list.Elems) < 6 {
		return nil, fmt.Errorf("obis tlv has only %d elements", len(list.Elems))
	}

	exact, simplified, err := parseKey(list.Elems[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse key: %w", err)
	}

	name := obis.GetFromCatalogue(exact)
	if name == "" {
		name = obis.GetFromCatalogue(simplified)
	}

	entry := &obis.OBISEntry{
		Name:          name,
		ExactKey:      exact,
		SimplifiedKey: simplified,
	}

	unit := parseUnit(list.Elems[3])

	scale := parseScale(list.Elems[4])

	err = parseValue(entry, unit, scale, list.Elems[5])
	if err != nil {
		return nil, fmt.Errorf("could not parse value: %w", err)
	}

	return entry, nil

}

func parseUnit(tlv *TLV) uint8 {
	if len(tlv.Value) < 1 {
		return 0
	}
	return uint8(tlv.Value[0])
}

func parseScale(tlv *TLV) int {
	if len(tlv.Value) < 1 {
		return 0
	}
	scale8 := int8(tlv.Value[0])
	return int(scale8)
}

func parseValue(entry *obis.OBISEntry, unitnum uint8, scale int, value *TLV) error {
	unit := GetUnitByKey(unitnum)
	err := unit.SetValue(entry, scale, value.Value, entry.SimplifiedKey)
	if err != nil {
		return err
	}
	return nil
}

func parseKey(key *TLV) (string, string, error) {

	if key.Type != TLVType_OctetStream {
		return "", "", errors.New("key tlv is not of type octet-stream")
	}

	v := key.Value
	if len(v) < 5 {
		return "", "", fmt.Errorf("obis tlv has only length of %d", len(v))
	}

	medium := uint8(v[0])
	channel := uint8(v[1])

	number1 := uint8(v[2])
	number2 := uint8(v[3])
	number3 := uint8(v[4])

	tariff := uint8(0)
	if len(v) >= 6 {
		tariff = uint8(v[5])
	}

	exactKey := fmt.Sprintf("%d-%d:%d.%d.%d*%d",
		medium, channel, number1, number2, number3, tariff)

	simplifiedKey := fmt.Sprintf("%d.%d.%d",
		number1, number2, number3)

	return exactKey, simplifiedKey, nil

}
