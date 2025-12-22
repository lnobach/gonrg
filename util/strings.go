package util

func BytesToPrintableString(b []byte) string {
	result := make([]byte, len(b))

	for i, v := range b {
		if v >= 32 && v <= 126 {
			result[i] = v
		} else {
			result[i] = '.'
		}
	}

	return string(result)
}
