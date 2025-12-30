package obis

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

	return &OBISMappedResult{DeviceID: list.DeviceID, MeasurementTime: list.MeasurementTime,
		List: list.List, Map: obismap}

}
