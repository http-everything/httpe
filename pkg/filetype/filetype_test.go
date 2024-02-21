package filetype_test

import (
	"testing"

	"http-everything/httpe/pkg/filetype"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiletypes(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{
			name: "image.jpg",
			want: "image/jpeg",
		},
		{
			name: "image.png",
			want: "image/png",
		},
		{
			name: "doc.pdf",
			want: "application/pdf",
		},
		{
			name: "archive.zip",
			want: "application/zip",
		},
		{
			name: "archive.tar.gz",
			want: "application/gzip",
		},
		{
			name: "text.txt",
			want: "text/UTF-8",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filetype.Type("../../testdata/files/" + tc.name)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
