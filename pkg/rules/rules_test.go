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
	rules := rules.New(logger)
	_, ruleErr := rules.Load("../../testdata/rules/good/all.yaml")
	require.NoError(t, ruleErr)
	valErr := rules.Validate()
	require.NoError(t, valErr)
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
			},
			expectedMarshalErr: false,
		},
		{
			name: "wrong-enums",
			expectedErrors: []string{
				"0.on.methods.1: 0.on.methods.1 must be one of the following: \"get\", \"post\", \"put\", \"delete\", \"options\"",
				"0.with.auth_hashing: 0.with.auth_hashing must be one of the following: \"sha256\", \"sha512\"",
			},
			expectedMarshalErr: false,
		},
		{
			name:               "wrong-timeout",
			expectedErrors:     []string{"cannot unmarshal !!str `bla` into float64"},
			expectedMarshalErr: true,
		},
		{
			name:               "broken-yaml",
			expectedErrors:     []string{"mapping values are not allowed in this context"},
			expectedMarshalErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, logFile := makeTestLogger(t)
			ru := rules.New(logger)
			_, ruleErr := ru.Load("../../testdata/rules/bad/" + tc.name + ".yaml")
			t.Logf("Err %s", ruleErr)
			valErr := ru.Validate()
			log, err := os.ReadFile(logFile)
			require.NoError(t, err)
			t.Logf("Logged: %s", string(log))

			if tc.expectedMarshalErr {
				for _, exp := range tc.expectedErrors {
					assert.ErrorContains(t, ruleErr, exp)
				}
			} else {
				assert.ErrorContains(t, valErr, "invalid rules file")
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
