package timespan

func GetFirst[T any](arr []T) (T, bool) {
	if len(arr) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	return arr[0], true
}
