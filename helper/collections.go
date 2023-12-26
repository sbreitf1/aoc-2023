package helper

import "sort"

func GetReversedSlice[T any](arr []T) []T {
	arr2 := make([]T, len(arr))
	copy(arr2, arr)
	ReverseSlice(arr2)
	return arr2
}

func ReverseSlice[T any](arr []T) {
	for i := 0; i < len(arr)/2; i++ {
		tmp := arr[i]
		arr[i] = arr[len(arr)-i-1]
		arr[len(arr)-i-1] = tmp
	}
}

func CloneMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func IterateMapInKeyOrder[K Ordered, V any](m map[K]V, f func(k K, v V)) {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for _, k := range keys {
		f(k, m[k])
	}
}

func LinesToRunes(lines []string) [][]rune {
	runeLines := make([][]rune, len(lines))
	for y := 0; y < len(lines); y++ {
		runeLines[y] = []rune(lines[y])
	}
	return runeLines
}

func RunesToLines(runeLines [][]rune) []string {
	lines := make([]string, len(runeLines))
	for y := 0; y < len(lines); y++ {
		lines[y] = string(runeLines[y])
	}
	return lines
}
