package timeunit

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a string duration with days, hours, minutes, etc.
func ParseDuration(s string) (time.Duration, error) {
	unitMultiplier := map[string]time.Duration{
		"s":   time.Second,
		"min": time.Minute,
		"h":   time.Hour,
		"d":   time.Hour * 24,
		"w":   time.Hour * 24 * 7,
		"m":   time.Hour*24*30 + time.Hour*12,
	}

	// Find the unit in the string
	for unit, multiplier := range unitMultiplier {
		if strings.HasPrefix(s, "-") {
			return 0, fmt.Errorf("%s: negative values are not suppotrted", s)
		}
		if strings.HasSuffix(s, unit) {
			valueStr := strings.TrimSuffix(s, unit)
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid duration value: %s", s)
			}
			return time.Duration(value) * multiplier, nil
		}
	}

	return 0, fmt.Errorf("%s: invalid duration format or unsupported unit", s)
}
