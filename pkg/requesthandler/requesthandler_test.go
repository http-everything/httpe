package requesthandler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/share/logger"

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
		{
			name: "Postaction",
			rule: rules.Rule{
				AnswerContent: "1",
				PostAction:    &rules.PostAction{RunScript: "echo 1"}},
			wantBody:   "1",
			wantStatus: http.StatusOK,
		},
	}
	require.NoError(t, err)
	dataDir := t.TempDir()
	logFile := filepath.Join(t.TempDir(), "test.log")
	l, err := logger.New("test", logFile, logger.DEBUG)
	require.NoError(t, err)
	conf := config.Config{S: &config.SvrConfig{DataDir: dataDir, DataRetention: "1s"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.rule.On = &rules.On{
				Path: "/",
			}
			req, err := http.NewRequest("get", "/", nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			httpHandler := requesthandler.Execute(tc.rule, l, &conf)
			httpHandler.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantBody, rec.Body.String())
			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
	// Look for files created by asynchronous post actions in the data dir
	if runtime.GOOS == "windows" {
		// Windows is slow
		time.Sleep(600 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	files, err := os.ReadDir(dataDir)
	require.NoError(t, err)

	// Because one test case contains one post action, we expect one file to be written
	assert.Len(t, files, 1, "no files in data directory "+dataDir)

	// Shut down the logger. Temp dir cannot be removed if the log is still open
	l.Shutdown()

	// Read the log file, it shall not contain errors
	log, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.NotContains(t, log, "ERROR")
}

func newline(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
