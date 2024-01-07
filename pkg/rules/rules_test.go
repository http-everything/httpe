package rules_test

import (
	"fmt"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldSucceed(t *testing.T) {
	logger, logFile := makeTestLogger(t)
	ru := rules.New(logger)
	_, ruleErr := ru.LoadAndValidate("../../testdata/rules/good/all.yaml")
	require.NoError(t, ruleErr)
	log, err := os.ReadFile(logFile)
	require.NoError(t, err)

	assert.Contains(t, string(log), "'../../testdata/rules/good/all.yaml' successfully validated against schema")
}

func TestShouldFail(t *testing.T) {
	cases := []struct {
		name               string
		expectedErrors     []string
		expectedMarshalErr bool
	}{
		{
			name:               "name-missing",
			expectedErrors:     []string{"name is required"},
			expectedMarshalErr: false,
		},
		{
			name: "everything-missing",
			expectedErrors: []string{
				"name is required",
				"on is required",
				"do is required",
			},
			expectedMarshalErr: false,
		},
		{
			name:               "wrong-timeout",
			expectedErrors:     []string{"cannot unmarshal !!str `bla` into float64"},
			expectedMarshalErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, logFile := makeTestLogger(t)
			ru := rules.New(logger)
			_, ruleErr := ru.LoadAndValidate("../../testdata/rules/bad/" + tc.name + ".yaml")
			t.Logf("Err %s", ruleErr)
			log, err := os.ReadFile(logFile)
			require.NoError(t, err)
			t.Logf("Logged: %s", string(log))

			if tc.expectedMarshalErr {
				for _, exp := range tc.expectedErrors {
					assert.ErrorContains(t, ruleErr, exp)
				}
			} else {
				assert.ErrorContains(t, ruleErr, "invalid rules file")
				assert.Contains(t, string(log), fmt.Sprintf("ERROR: test: schema validation against %s failed", rules.SchemaURL))
				for _, exp := range tc.expectedErrors {
					assert.Contains(t, string(log), exp)
				}
			}
		})
	}
}

func makeTestLogger(t *testing.T) (l *logger.Logger, logFile string) {
	t.Helper()
	logFile = t.TempDir() + "/test.log"

	l, err := logger.New("test", logFile, "debug")
	require.NoError(t, err)
	return l, logFile
}
