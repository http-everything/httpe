package requesthandler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"

	"github.com/http-everything/httpe/pkg/requesthandler"
	"github.com/http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestHandler(t *testing.T) {
	dummyFile := t.TempDir() + "/dummy.txt"
	err := os.WriteFile(dummyFile, []byte("test"), 0400)
	require.NoError(t, err)

	cases := []struct {
		name       string
		rule       rules.Rule
		wantBody   string
		wantStatus int
	}{
		{
			name:       "Answer Content",
			rule:       rules.Rule{AnswerContent: "foo"},
			wantBody:   "foo",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Answer File",
			rule:       rules.Rule{AnswerFile: dummyFile},
			wantBody:   "test",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Redir Perm",
			rule:       rules.Rule{RedirectPermanent: "/test"},
			wantBody:   "",
			wantStatus: http.StatusMovedPermanently,
		},
		{
			name:       "Redir Temp",
			rule:       rules.Rule{RedirectTemporary: "/test"},
			wantBody:   "",
			wantStatus: http.StatusFound,
		},
		{
			name:       "Run Script",
			rule:       rules.Rule{RunScript: "echo test"},
			wantBody:   "test" + newline(t),
			wantStatus: http.StatusOK,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.rule.On = &rules.On{
				Path: "/",
			}
			//rule := rules.Rule{
			//	On: &rules.On{
			//		Path: "/",
			//	},
			//	AnswerContent: "foo",
			//}

			req, err := http.NewRequest("get", "/", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			httpHandler := requesthandler.Execute(tc.rule, nil, nil)
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
