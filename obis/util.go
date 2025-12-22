package obis

import (
	"math"
)

func Floatify(e *OBISEntry) {
	e.ValueFloat = float64(e.ValueNum) * math.Pow10(e.ValueScale)
}
