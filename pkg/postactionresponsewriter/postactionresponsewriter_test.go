package postactionresponsewriter_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/postactionresponsewriter"
	"github.com/http-everything/httpe/pkg/share/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPostActionResponseWriter(t *testing.T) {
	logFile := t.TempDir() + "/test.log"
	l, err := logger.New("test", logFile, logger.DEBUG)
	require.NoError(t, err)
	dataDir := t.TempDir()
	conf := config.Config{S: &config.SvrConfig{DataDir: dataDir, DataRetention: "1s"}}

	for try, _ := range []int{1, 2} {
		t.Run(fmt.Sprintf("Round %d", try), func(t *testing.T) {
			writer := postactionresponsewriter.New(&conf, l)
			writer.AddActionResponse("test", actions.ActionResponse{
				SuccessBody: "1",
				Code:        0,
				ErrorBody:   "2",
			}, errors.New("3"))
			writer.Write()
			logBytes, err := os.ReadFile(logFile)
			require.NoError(t, err)
			log := string(logBytes)

			t.Logf("Content of log file %s:\n%s", logFile, log)
			assert.NotContains(t, log, "ERROR")

			files, err := os.ReadDir(dataDir)
			require.NoError(t, err)

			var parFile string
			for _, file := range files {
				if filepath.Ext(file.Name()) == ".json" {
					parFile = file.Name()
					break
				}
			}
			require.NotEmpty(t, parFile)
			t.Log(parFile)
			parBytes, err := os.ReadFile(dataDir + "/" + parFile)
			require.NoError(t, err)
			par := make([]postactionresponsewriter.PostActionResponse, 0)
			err = json.Unmarshal(parBytes, &par)
			require.NoError(t, err, "error unmarshalling file %s", parFile)

			// Because of auto-cleanup the second round will delete the file from the first round
			assert.Len(t, files, 1, "number of files in data dir")
			assert.Contains(t, log, fmt.Sprintf("Deleted %d files from data directory", try))
			assert.Equal(t, par[0].SuccessBody, "1")
			assert.Equal(t, par[0].Code, 0)
			assert.Equal(t, par[0].ErrorBody, "2")
			assert.Equal(t, par[0].InternalError, "3")
			if try == 0 {
				time.Sleep(2 * time.Second)
			}
		})
	}
	// Shut down the logger. Temp dir cannot be removed if the log is still open
	l.Shutdown()
}
