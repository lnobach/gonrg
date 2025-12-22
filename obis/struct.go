package obis

import (
	"fmt"
	"time"
)

type OBISEntry struct {
	ExactKey      string  `json:"exactKey"`
	SimplifiedKey string  `json:"simplifiedKey"`
	Name          string  `json:"name"`
	ValueText     string  `json:"valueText"`
	ValueNum      int64   `json:"valueNum"`
	ValueScale    int     `json:"valueScale"`
	ValueFloat    float64 `json:"valueFloat"`
	Unit          string  `json:"unit"`
}

type OBISMap map[string]*OBISEntry

type OBISList []*OBISEntry

type OBISListResult struct {
	MeasurementTime time.Time `json:"measurementTime"`
	DeviceID        string    `json:"deviceID"`
	List            OBISList  `json:"list"`
}

type OBISMappedResult struct {
	MeasurementTime time.Time `json:"measurementTime"`
	DeviceID        string    `json:"deviceID"`
	List            OBISList  `json:"list"`
	Map             OBISMap   `json:"map"`
}

type OBISSingleResult struct {
	MeasurementTime time.Time  `json:"measurementTime"`
	DeviceID        string     `json:"deviceID"`
	Entry           *OBISEntry `json:"list"`
}

func (e *OBISEntry) PrettyValue(unit bool) string {
	if e.ValueText != "" {
		return e.ValueText
	}
	if e.ValueNum == 0 && e.Unit == "" {
		return "-"
	}
	num := fmt.Sprintf("%d", e.ValueNum)

	for {
		//nolint:staticcheck // won't run into loop
		if e.ValueScale < len(num) {
			break
		}
		num = "0" + num
	}

	if e.ValueScale < 0 {
		scaleRev := len(num) + e.ValueScale
		num = num[:scaleRev] + "." + num[scaleRev:]
	}
	if unit && e.Unit != "" {
		return num + " " + e.Unit
	}
	return num
}

func (m *OBISMappedResult) GetList() *OBISListResult {
	return &OBISListResult{
		DeviceID:        m.DeviceID,
		List:            m.List,
		MeasurementTime: m.MeasurementTime,
	}
}
