package postaction_test

import (
	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/postaction"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunScript(t *testing.T) {
	r := rules.Rule{
		PostAction: &rules.PostAction{RunScript: "echo 1"},
	}
	rd := requestdata.Data{}
	dataDir := t.TempDir()
	logFile := filepath.Join(t.TempDir(), "log.txt")
	l, err := logger.New("test", logFile, logger.DEBUG)
	require.NoError(t, err)
	conf := config.Config{S: &config.SvrConfig{DataDir: dataDir, DataRetention: "1s"}}
	postaction.Execute(r, rd, &conf, l)
	l.Shutdown()
	// Look for files created by asynchronous post actions in the data dir
	time.Sleep(100 * time.Millisecond)
	files, err := os.ReadDir(dataDir)
	require.NoError(t, err)

	// Because the test case contains one post action, we expect one file to be written
	assert.Len(t, files, 1, "no files in data directory")

	// Read the log file, it shall not contain errors
	log, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.NotContains(t, log, "ERROR")
}
