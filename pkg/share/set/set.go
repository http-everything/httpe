package set

// String is a function that returns the first non-empty string from the provided list of strings.
// If all values are empty, it returns an empty string.
func String(values ...string) string {
	for _, val := range values {
		if val != "" {
			return val
		}
	}
	return ""
}

// Int is a function that returns the first non-zero integer from the provided list of integers.
// If all values are zero, it returns zero.
func Int(values ...int) int {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}
