package helper

func GreatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LeastCommonMultiple(vals ...int64) int64 {
	result := vals[0] * vals[1] / GreatestCommonDivisor(vals[0], vals[1])
	for i := 2; i < len(vals); i++ {
		result = LeastCommonMultiple(result, vals[i])
	}
	return result
}

func Min[T Ordered](values ...T) T {
	min := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
	}
	return min
}

func Max[T Ordered](values ...T) T {
	max := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}
	return max
}
