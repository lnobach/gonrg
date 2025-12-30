package obis

import "maps"

func ListToMap(list *OBISListResult) *OBISMappedResult {

	obismap_exact := make(OBISMap)
	obismap := make(OBISMap)

	for _, e := range list.List {
		obismap_exact[e.ExactKey] = e
		obismap[e.SimplifiedKey] = e
		if e.Name != "" {
			obismap[e.Name] = e
		}
	}

	maps.Copy(obismap, obismap_exact)

	return &OBISMappedResult{DeviceID: list.DeviceID, MeasurementTime: list.MeasurementTime,
		List: list.List, Map: obismap}

}
