package util

import "strings"

func KeysFromMap[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func KeysToStr[T any](m map[string]T, sep string) string {
	return strings.Join(KeysFromMap(m), sep)
}
