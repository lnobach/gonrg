package obis

// Returns nil if nothing has changed, or a list result
// with the changed elements otherwise.
func GetChanged(previous, next *OBISListResult) *OBISListResult {

	changed := make([]*OBISEntry, 0, 20)

	if previous == nil {
		return next
	}

	for i, e := range next.List {
		if i >= len(previous.List) {
			changed = append(changed, e)
			continue
		}
		if HasChanged(previous.List[i], e) {
			changed = append(changed, e)
		}
	}

	if len(changed) <= 0 {
		return nil
	}

	return &OBISListResult{
		MeasurementTime: next.MeasurementTime,
		DeviceID:        next.DeviceID,
		List:            changed,
	}

}

func HasChanged(previous, next *OBISEntry) bool {
	if previous.ValueNum != next.ValueNum {
		return true
	}
	if previous.ValueText != next.ValueText {
		return true
	}
	if previous.ValueScale != next.ValueScale {
		return true
	}
	if previous.ExactKey != next.ExactKey {
		return true
	}
	if previous.Unit != next.Unit {
		return true
	}
	return false
}
