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

// Like strings.EqualFold, but only non-unicode i.e. ascii
// characters are folded
func EqualFoldNonUnicode(s, t string) bool {
	i := 0
	for n := min(len(s), len(t)); i < n; i++ {
		sr := s[i]
		tr := t[i]

		if tr == sr {
			continue
		}

		if tr < sr {
			tr, sr = sr, tr
		}

		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			continue
		}
		return false
	}
	return len(s) == len(t)
}
