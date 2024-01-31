package logger_test

import (
	"fmt"
	"http-everything/httpe/pkg/share/logger"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerInfoLevel(t *testing.T) {
	cases := []string{
		"info",
		"debug",
		"error",
	}
	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			logFile := t.TempDir() + "/httpe.log"
			l, err := logger.New("test", logFile, tc)
			assert.NoError(t, err)
			l.Infof("--info--")
			l.Errorf("--error--")
			l.Debugf("--debug--")
			lf, err := os.ReadFile(logFile)
			assert.NoError(t, err)

			require.FileExists(t, logFile)
			require.Contains(t, string(lf), fmt.Sprintf("%s: test: --%s--", strings.ToUpper(tc), tc))
			l.Shutdown()
		})
	}
}

func TestLoggerInfoError(t *testing.T) {
	logFile := t.TempDir() + "/httpe.log"
	l, err := logger.New("test", logFile, "error")
	require.NoError(t, err)
	l.Infof("1234abc")
	lf, err := os.ReadFile(logFile)

	assert.NoError(t, err)
	assert.Empty(t, string(lf))
	l.Shutdown()
}
