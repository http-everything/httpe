package remove_test

import (
	"testing"

	"github.com/http-everything/httpe/pkg/share/remove"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFileExtension(t *testing.T) {
	cases := []struct {
		input     string
		wants     string
		extension string
	}{
		{
			input:     "powershell.exe",
			wants:     "powershell",
			extension: "exe",
		},
		{
			input:     "SomeFile",
			wants:     "SomeFile",
			extension: "jpg",
		},
		{
			input:     "my.file.pNg",
			wants:     "my.file",
			extension: "PnG",
		},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			out := remove.FileExtension(tc.input, tc.extension)

			assert.Equal(t, tc.wants, out)
		})
	}
}
