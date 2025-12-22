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

func AllCharsPrintable(b []byte) bool {
	for _, v := range b {
		if v < 32 || v > 126 {
			return false
		}
	}
	return true
}
