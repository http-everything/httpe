package extract

import (
	"fmt"
	"reflect"
	"strings"
)

// SFromI extracts string values from unstructured interfaces using dot notation.
// Example: SFromI("address.city",interface)
func SFromI(path string, input interface{}) string {
	// Split the path
	keys := strings.Split(path, ".")

	// Start with the input data
	var data interface{} = input

	// Traverse the keys
	for _, key := range keys {
		// Handle nil data
		if data == nil {
			return ""
		}

		// Get the reflected Value
		val := reflect.ValueOf(data)

		// Handle invalid types
		if val.Kind() != reflect.Map {
			return ""
		}

		// Lookup the key
		value := val.MapIndex(reflect.ValueOf(key))

		if !value.IsValid() {
			return ""
		}

		// Save value for next lookup
		data = value.Interface()
	}

	// Return final value as string
	return fmt.Sprintf("%v", data)
}
