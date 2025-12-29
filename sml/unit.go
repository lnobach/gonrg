package sml

import (
	"fmt"

	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/util"
)

var (
	units []SMLUnit
)

func initMap() {
	units = make([]SMLUnit, 256)

	units[0x1e] = &unit_uint64_1000{unit_uint64{name: "kWh"}}
	units[0x1b] = &unit_int64{name: "W"}
}

type SMLUnit interface {
	SetValue(e *obis.OBISEntry, scale int, raw []byte, simplifiedKey string) error
}

func GetDefaultUnit() SMLUnit {
	return &defaultunit{}
}

func GetUnitByKey(key byte) SMLUnit {
	if len(units) == 0 {
		initMap()
	}
	unit := units[key]
	if unit == nil {
		return GetDefaultUnit()
	}
	return unit
}

// === Fallback unit if unit is unknown
type defaultunit struct {
	name string
}

func (d *defaultunit) SetValue(e *obis.OBISEntry, scale int, raw []byte,
	simplifiedKey string) error {
	if isDeviceID(simplifiedKey) {
		id, err := parseDeviceID(raw)
		if err != nil {
			return fmt.Errorf("cannot parse device id: %w", err)
		}
		e.ValueText = id
		e.Unit = d.name
		return nil
	}
	e.ValueText = fmt.Sprintf("%x", raw)
	e.Unit = d.name
	return nil
}

func isDeviceID(sk string) bool {
	return sk == "96.1.0" || sk == "0.0.9"
}

// === Unsigned Integer Unit
type unit_uint64 struct {
	name string
}

func (d *unit_uint64) SetValue(e *obis.OBISEntry, scale int, raw []byte, _ string) error {
	e.ValueNum = int64(util.BytesToUint64(raw))
	e.ValueScale = scale
	e.Unit = d.name
	return nil
}

// === Signed Integer Unit
type unit_int64 struct {
	name string
}

func (d *unit_int64) SetValue(e *obis.OBISEntry, scale int, raw []byte, _ string) error {
	e.ValueNum = util.BytesToInt64(raw)
	e.ValueScale = scale
	e.Unit = d.name
	return nil
}

// === Unsigned Integer Unit / 1000
type unit_uint64_1000 struct {
	unit_uint64
}

func (d *unit_uint64_1000) SetValue(e *obis.OBISEntry, scale int, raw []byte,
	simplifiedKey string) error {
	err := d.unit_uint64.SetValue(e, scale, raw, simplifiedKey)
	e.ValueScale = e.ValueScale - 3
	return err
}
