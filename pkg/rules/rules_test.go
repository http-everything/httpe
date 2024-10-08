package rules_test

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/http-everything/httpe/pkg/config"

	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/extract"
	"github.com/http-everything/httpe/pkg/share/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var smtpConfig = &config.SMTPConfig{
	Server: "127.0.0.1:25",
}

func TestRulesShouldSucceed(t *testing.T) {
	logger, logFile := makeTestLogger(t)
	files, err := os.ReadDir("../../testdata/rules/good/")

	require.NoError(t, err)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			rulesCfg, ruleErr := rules.Read("../../testdata/rules/good/"+file.Name(), logger)
			assert.NoError(t, ruleErr)
			valErr := rulesCfg.Validate(smtpConfig)
			assert.NoError(t, valErr)
			log, err := os.ReadFile(logFile)
			assert.NoError(t, err)

			assert.Contains(t, string(log), "successfully validated against schema")
		})
	}
	logger.Shutdown()
}

func TestRulesShouldFail(t *testing.T) {
	cases := []struct {
		name                    string
		wantErrors              []string
		wantMarshalError        bool
		wantJSONValidationError bool
	}{
		{
			name: "action-missing",
			wantErrors: []string{
				"rule 0 'Action missing' is missing a valid action.",
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
		{
			name: "missing-email-fields",
			wantErrors: []string{
				"to is required",
				"body is required",
			},
		},
		{
			name: "wrong-postaction",
			wantErrors: []string{
				"rule 0 'Wrong Postaction' invalid postaction",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger, logFile := makeTestLogger(t)
			rulesCfg, ruleErr := rules.Read("../../testdata/rules/bad/"+tc.name+".yaml", logger)
			var valErr error
			if ruleErr == nil {
				valErr = rulesCfg.Validate(smtpConfig)
			}

			log, err := os.ReadFile(logFile)
			require.NoError(t, err)
			logger.Shutdown()

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

func TestMissingSmtpConfig(t *testing.T) {
	logger, logFile := makeTestLogger(t)
	rulesCfg, ruleErr := rules.Read("../../testdata/rules/good/send-email.yaml", logger)
	assert.NoError(t, ruleErr)
	valErr := rulesCfg.Validate(nil)
	assert.ErrorContains(t, valErr, "invalid rules")
	log, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	assert.Contains(t, string(log), "send.email requires an smtp configuration in httpe configuration file")
	logger.Shutdown()
}

func TestYamlToJSON(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		JSON := rules.YamlToJSON("test.yaml")

		if runtime.GOOS == "windows" {
			assert.Equal(t, "open test.yaml: The system cannot find the file specified.", JSON)
		} else {
			assert.Equal(t, "open test.yaml: no such file or directory", JSON)
		}
	})

	t.Run("good yaml", func(t *testing.T) {
		yaml := t.TempDir() + "/test.yaml"
		err := os.WriteFile(yaml, []byte(`foo: test`), 0400)
		require.NoError(t, err)
		JSON := rules.YamlToJSON(yaml)
		var data interface{}
		err = json.Unmarshal([]byte(JSON), &data)
		require.NoError(t, err)

		assert.Equal(t, "test", extract.SFromI("foo", data))
	})

	t.Run("bad yaml", func(t *testing.T) {
		yaml := t.TempDir() + "/test.yaml"
		err := os.WriteFile(yaml, []byte(`&*`), 0400)
		require.NoError(t, err)
		JSON := rules.YamlToJSON(yaml)

		assert.Equal(t, "yaml: did not find expected alphabetic or numeric character", JSON)
	})
}

func makeTestLogger(t *testing.T) (l *logger.Logger, logFile string) {
	t.Helper()
	logFile = t.TempDir() + "/test.log"

	l, err := logger.New("test", logFile, "debug")
	require.NoError(t, err)
	return l, logFile
}
