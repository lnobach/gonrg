package obis

import (
	"encoding/json"
	"time"

	"github.com/lnobach/gonrg/util"
)

type rawJsonObj struct {
	formatFunc func() ([]byte, error)
}

func (f *rawJsonObj) MarshalJSON() ([]byte, error) {
	return f.formatFunc()
}

type OBISEntry struct {
	ExactKey      string `json:"exactKey"`
	SimplifiedKey string `json:"simplifiedKey"`
	Name          string `json:"name"`
	ValueText     string `json:"valueText"`
	ValueNum      int64  `json:"valueNum"`
	ValueScale    int    `json:"valueScale"`
	//	ValueFloat    float64 `json:"valueFloat"`
	Unit string `json:"unit"`
}

func (e *OBISEntry) MarshalJSON() ([]byte, error) {
	type Alias OBISEntry
	return json.Marshal(&struct {
		*Alias
		ValueFloat rawJsonObj `json:"valueFloat"`
	}{
		Alias: (*Alias)(e),
		ValueFloat: rawJsonObj{
			formatFunc: func() ([]byte, error) {
				return []byte(util.DecimalScaleToString(e.ValueNum, e.ValueScale)), nil
			},
		},
	})
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

	num := util.DecimalScaleToString(e.ValueNum, e.ValueScale)

	if unit && e.Unit != "" {
		return num + " " + e.Unit
	}
	return num
}

func (m *OBISMappedResult) GetList() *OBISListResult {
	if m == nil {
		return nil
	}
	return &OBISListResult{
		DeviceID:        m.DeviceID,
		List:            m.List,
		MeasurementTime: m.MeasurementTime,
	}
}
