package rules_test

import (
	"fmt"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldSucceed(t *testing.T) {
	logger, logFile := makeTestLogger(t)
	files, err := os.ReadDir("../../testdata/rules/good/")
	require.NoError(t, err)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		t.Run(file.Name(), func(T *testing.T) {
			rulesCfg, ruleErr := rules.Read("../../testdata/rules/good/"+file.Name(), logger)
			assert.NoError(t, ruleErr)
			valErr := rulesCfg.Validate()
			assert.NoError(t, valErr)
			log, err := os.ReadFile(logFile)
			assert.NoError(t, err)

			assert.Contains(t, string(log), "successfully validated against schema")
		})
	}
}

func TestShouldFail(t *testing.T) {
	cases := []struct {
		name                    string
		wantErrors              []string
		wantMarshalError        bool
		wantJSONValidationError bool
	}{
		{
			name: "action-missing",
			wantErrors: []string{
				"rule 0 'Action missing' is missing a valid action in the 'do' section.",
				fmt.Sprintf("Use one of '%s'", strings.Join(rules.ValidActions, ", ")),
			},
			wantMarshalError:        false,
			wantJSONValidationError: false,
		},
		{
			name: "wrong-values",
			wantErrors: []string{
				"0.on.methods.1: rules.0.on.methods.1 must be one of the following: \"get\", \"post\", \"put\", \"delete\", \"options\"",
				"0.with.auth_hashing: rules.0.with.auth_hashing must be one of the following: \"sha256\", \"sha512\"",
			},
			wantMarshalError: false,
		},
		{
			name:             "wrong-timeout",
			wantErrors:       []string{"cannot unmarshal !!str `bla` into int"},
			wantMarshalError: true,
		},
		{
			name:             "broken-yaml",
			wantErrors:       []string{"mapping values are not allowed in this context"},
			wantMarshalError: true,
		},
		{
			name: "wrong-headers",
			wantErrors: []string{
				"Additional property White Space is not allowed",
				"Additional property Föö is not allowed",
			},
		},
		{
			name: "wrong-bytes",
			wantErrors: []string{
				"max_request_body: Does not match pattern '^[0-9]+ ?[BKMGTP]{0,2}$'\n",
			},
		},
		{
			name: "wrong-anchor",
			wantErrors: []string{
				"yaml: unknown anchor 'auth' referenced",
			},
			wantMarshalError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, logFile := makeTestLogger(t)
			rulesCfg, ruleErr := rules.Read("../../testdata/rules/bad/"+tc.name+".yaml", logger)
			var valErr error
			if ruleErr == nil {
				valErr = rulesCfg.Validate()
			}

			log, err := os.ReadFile(logFile)
			require.NoError(t, err)

			if tc.wantMarshalError {
				for _, exp := range tc.wantErrors {
					assert.ErrorContains(t, ruleErr, exp)
				}
			} else {
				assert.ErrorContains(t, valErr, "invalid rules")
				for _, exp := range tc.wantErrors {
					assert.Contains(t, string(log), exp)
				}
			}
			if tc.wantJSONValidationError {
				assert.Contains(t, string(log), fmt.Sprintf("ERROR: test: schema validation against %s failed", rules.SchemaURL))
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
