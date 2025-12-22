package sml

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/lnobach/gonrg/obis"
)

var (
	units []SMLUnit
)

func initMap() {
	units = make([]SMLUnit, 256)

	units[0x1e] = &unit_uint32_1000{unit_uint32{name: "kWh"}}
	units[0x1b] = &unit_uint32{name: "W"}
}

type SMLUnit interface {
	AsString() string
	SetValue(e *obis.OBISEntry, scale int, raw []byte) error
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

func (d *defaultunit) AsString() string {
	return d.name
}

func (d *defaultunit) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	e.ValueText = fmt.Sprintf("%x", raw)
	return nil
}

// === Integer Unit
type unit_uint32 struct {
	name string
}

func (d *unit_uint32) AsString() string {
	return d.name
}

func (d *unit_uint32) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	if len(raw) < 8 {
		raw = append(bytes.Repeat([]byte{0x00}, 8-len(raw)), raw...)
	}
	e.ValueNum = int64(binary.BigEndian.Uint64(raw))
	e.ValueScale = scale
	return nil
}

// === Integer Unit / 1000
type unit_uint32_1000 struct {
	unit_uint32
}

func (d *unit_uint32_1000) AsString() string {
	return d.name
}

func (d *unit_uint32_1000) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	err := d.unit_uint32.SetValue(e, scale, raw)
	e.ValueScale = e.ValueScale - 3
	return err
}
