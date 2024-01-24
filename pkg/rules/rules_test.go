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
	rules := rules.New(logger)
	_, ruleErr := rules.Load("../../testdata/rules/good/hello-world.yaml")
	require.NoError(t, ruleErr)
	valErr := rules.Validate()
	require.NoError(t, valErr)
	log, err := os.ReadFile(logFile)
	require.NoError(t, err)

	assert.Contains(t, string(log), "'../../testdata/rules/good/hello-world.yaml' successfully validated against schema")
}

func TestShouldFail(t *testing.T) {
	cases := []struct {
		name                    string
		wantErrors              []string
		wantMarshalError        bool
		wantJSONValidationError bool
	}{
		{
			name:             "name-missing",
			wantErrors:       []string{"name is required"},
			wantMarshalError: false,
		},
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
			name: "wrong-enums",
			wantErrors: []string{
				"0.on.methods.1: 0.on.methods.1 must be one of the following: \"get\", \"post\", \"put\", \"delete\", \"options\"",
				"0.with.auth_hashing: 0.with.auth_hashing must be one of the following: \"sha256\", \"sha512\"",
			},
			wantMarshalError: false,
		},
		{
			name: "wrong-auth-hashing",
			wantErrors: []string{
				"0.with.auth_hashing must be one of the following: \"sha256\", \"sha512\"",
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

			if tc.wantMarshalError {
				for _, exp := range tc.wantErrors {
					assert.ErrorContains(t, ruleErr, exp)
				}
			} else {
				assert.ErrorContains(t, valErr, "invalid rules file")
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
