package merge

import "strings"

// StringMapsI merges a list of maps into a single map where all keys are treated case-insensitive,
// keys in the returned result are converted to lower case
func StringMapsI(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			k = strings.ToLower(k)
			if _, ok := result[k]; !ok {
				result[k] = v
			}
		}
	}
	return result
}
