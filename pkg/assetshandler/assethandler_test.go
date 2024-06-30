package assetshandler_test

import (
	_ "embed"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/http-everything/httpe/pkg/assetshandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed alpine3.min.js
var alpineJS string

//go:embed bootstrap5.2.3.min.css
var bootstrapCSS string

//go:embed favicon.ico
var favicon string

//go:embed bootstrap5.2.3.bundle.min.js
var boostrapBundleJS string

func TestAssetsHandler(t *testing.T) {
	cases := []struct {
		requestURI  string
		wantStatus  int
		wantContent string
		wantHeader  string
	}{
		{
			requestURI:  "/_assets/alpine.js",
			wantStatus:  http.StatusOK,
			wantContent: alpineJS,
			wantHeader:  "application/javascript",
		},
		{
			requestURI:  "/_assets/bootstrap.css",
			wantStatus:  http.StatusOK,
			wantContent: bootstrapCSS,
			wantHeader:  "text/css",
		},
		{
			requestURI:  "/_assets/bootstrap.bundle.js",
			wantStatus:  http.StatusOK,
			wantContent: boostrapBundleJS,
			wantHeader:  "application/javascript",
		},
		{
			requestURI:  "/favicon.ico",
			wantStatus:  http.StatusOK,
			wantContent: favicon,
			wantHeader:  "image/x-icon",
		},
		{
			requestURI:  "/notfound",
			wantStatus:  http.StatusNotFound,
			wantContent: "/notfound not found\n",
			wantHeader:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.requestURI, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.requestURI, nil)
			rr := httptest.NewRecorder()

			assetshandler.AssetsHandler(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
			assert.Equal(t, tc.wantContent, string(body))
			if tc.wantHeader != "" {
				assert.Equal(t, tc.wantHeader, resp.Header.Get("Content-Type"))
			}
		})
	}
}
