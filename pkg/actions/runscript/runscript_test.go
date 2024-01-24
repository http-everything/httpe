package runscript_test

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/runscript"
	"http-everything/httpe/pkg/rules"
	"testing"

	"github.com/stretchr/testify/assert"

	"http-everything/httpe/pkg/requestdata"

	"github.com/stretchr/testify/require"
)

func TestScript_Execute(t *testing.T) {
	cases := []struct {
		name            string
		script          string
		timeout         int
		interpreter     string
		wantSuccessBody string
		wantErrorBody   string
		wantExitCode    int
		wantError       string
	}{
		{
			name:            "Bash Script succeeded",
			script:          `echo "Hello, world"`,
			timeout:         1,
			interpreter:     "/bin/bash",
			wantSuccessBody: "Hello, world\n",
			wantExitCode:    0,
		},
		{
			name:            "Python Script succeeded",
			script:          `print("Hello, world")`,
			timeout:         1,
			interpreter:     "python3",
			wantSuccessBody: "Hello, world\n",
			wantExitCode:    0,
		},
		{
			name:          "Bash Script failed",
			script:        `nonsense`,
			timeout:       1,
			interpreter:   "/bin/bash",
			wantErrorBody: "/bin/bash: line 1: nonsense: command not found\n",
			wantExitCode:  127,
		},
		{
			name:         "Sh Script timed out",
			script:       `sleep 2`,
			timeout:      1,
			wantExitCode: -1,
			wantError:    "script killed, timeout 1 sec exceeded",
		},
		{
			name:         "Sh Script timed out and error",
			script:       `nonsense;sleep 2`,
			timeout:      1,
			wantExitCode: -1,
			wantError:    "script killed, timeout 1 sec exceeded, sh: line 1: nonsense: command not found\n",
		},
	}
	reqData := requestdata.Data{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule := rules.Rule{
				Do: &rules.Do{
					RunScript: tc.script,
					Args: rules.Args{
						Interpreter: tc.interpreter,
						Timeout:     tc.timeout,
					},
				},
			}

			//Create the actioner that implements the action interface
			var actioner actions.Actioner = runscript.Script{}

			// Execute the action by calling the mandatory function Execute()
			actionResp, err := actioner.Execute(rule, reqData)

			if tc.wantError == "" {
				require.NoError(t, err)
				assert.Equal(t, actions.ActionResponse{
					SuccessBody: tc.wantSuccessBody,
					ErrorBody:   tc.wantErrorBody,
					Code:        tc.wantExitCode,
				}, actionResp)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
		})
	}
}
