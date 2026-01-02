package obis

import (
	"math"
)

func GetValueFloat(e *OBISEntry) float64 {
	return float64(e.ValueNum) * math.Pow10(e.ValueScale)
}
