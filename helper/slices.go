package helper

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
