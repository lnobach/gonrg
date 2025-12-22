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

	units[0x1e] = &unit_uint32_1000{unit_uint32{name: "kWh"}}
	units[0x1b] = &unit_uint32{name: "W"}
}

type SMLUnit interface {
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

func (d *defaultunit) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	e.ValueText = fmt.Sprintf("%x", raw)
	e.Unit = d.name
	return nil
}

// === Integer Unit
type unit_uint32 struct {
	name string
}

func (d *unit_uint32) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	e.ValueNum = int64(util.Uint64FromVarByteArray(raw))
	e.ValueScale = scale
	e.Unit = d.name
	return nil
}

// === Integer Unit / 1000
type unit_uint32_1000 struct {
	unit_uint32
}

func (d *unit_uint32_1000) SetValue(e *obis.OBISEntry, scale int, raw []byte) error {
	err := d.unit_uint32.SetValue(e, scale, raw)
	e.ValueScale = e.ValueScale - 3
	return err
}
