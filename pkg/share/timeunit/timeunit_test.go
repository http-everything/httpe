package timeunit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseDuration tests the ParseDuration function with various inputs
func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		err      bool
	}{
		{"1s", time.Second, false},
		{"1min", time.Minute, false},
		{"1h", time.Hour, false},
		{"1d", time.Hour * 24, false},
		{"1w", time.Hour * 24 * 7, false},
		{"1m", time.Hour*24*30 + time.Hour*12, false}, // 30.5 days
		{"2m", 2 * (time.Hour*24*30 + time.Hour*12), false},
		{"10", 10 * time.Minute, true},
		{"", 0, true},      // empty input
		{"1x", 0, true},    // invalid unit
		{"-1d", 0, true},   // invalid negative duration
		{"1.5h", 0, true},  // float value not supported
		{"1hour", 0, true}, // unsupported text format
		{"ou", 0, true},    // unsupported format
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseDuration(tt.input)

			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
