package requesthandler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"

	"http-everything/httpe/pkg/requesthandler"
	"http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestHandler(t *testing.T) {
	dummyFile := t.TempDir() + "/dummy.txt"
	err := os.WriteFile(dummyFile, []byte("test"), 0400)
	require.NoError(t, err)

	cases := []struct {
		name       string
		do         *rules.Do
		wantBody   string
		wantStatus int
	}{
		{
			name: "Answer Content",
			do: &rules.Do{
				AnswerContent: "foo",
			},
			wantBody:   "foo",
			wantStatus: http.StatusOK,
		},
		{
			name: "Answer File",
			do: &rules.Do{
				AnswerFile: dummyFile,
			},
			wantBody:   "test",
			wantStatus: http.StatusOK,
		},
		{
			name: "Redir Perm",
			do: &rules.Do{
				RedirectPermanent: "/test",
			},
			wantBody:   "",
			wantStatus: http.StatusMovedPermanently,
		},
		{
			name: "Redir Temp",
			do: &rules.Do{
				RedirectTemporary: "/test",
			},
			wantBody:   "",
			wantStatus: http.StatusFound,
		},
		{
			name: "Run Script",
			do: &rules.Do{
				RunScript: "echo test",
			},
			wantBody:   "test" + newline(t),
			wantStatus: http.StatusOK,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule := rules.Rule{
				On: &rules.On{
					Path: "/",
				},
				Do: tc.do,
			}

			req, err := http.NewRequest("get", "/", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			httpHandler := requesthandler.Execute(rule, nil)
			httpHandler.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantBody, rec.Body.String())
			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func newline(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
