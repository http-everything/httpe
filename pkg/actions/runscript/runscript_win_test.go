//go:build windows

package runscript_test

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/runscript"
	"http-everything/httpe/pkg/rules"
	"testing"

	"github.com/stretchr/testify/assert"

	"http-everything/httpe/pkg/requestdata"
)

const PowerShell = "powershell"

func TestScriptExecuteWindows(t *testing.T) {
	var dir = "C:\\Windows\\temp"

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
			name:            "PowerShell Script succeeded",
			script:          `Write-Output "Hello, world"`,
			timeout:         1,
			interpreter:     PowerShell,
			wantSuccessBody: "Hello, world\r\n",
			wantExitCode:    0,
		},
		{
			name:            "PowerShell Script succeeded (multi line)",
			script:          "Write-Output \"Line1\"\nWrite-Output \"Line2\"",
			timeout:         1,
			interpreter:     PowerShell,
			wantSuccessBody: "Line1\r\nLine2\r\n",
			wantExitCode:    0,
		},
		{
			name:          "PowerShell Script failed",
			script:        "gaga",
			timeout:       1,
			wantErrorBody: "The term 'gaga' is not recognized as the name of a cmdlet",
			wantExitCode:  1,
		},
		{
			name:         "Powershell Script timed out",
			script:       `Start-Sleep 2`,
			timeout:      1,
			wantExitCode: -1,
			wantError:    "process error: exit status 1",
		},
		{
			name:            "Verify directory",
			script:          "(Get-Location).Path",
			wantSuccessBody: "C:\\Windows\\temp\r\n",
		},
		{
			name:            "Cmd.exe Script succeeded",
			script:          "echo hello world",
			timeout:         1,
			interpreter:     "cmd.exe",
			wantSuccessBody: "hello world\r\n",
			wantExitCode:    0,
		},
		{
			name:            "Cmd Script succeeded (no suffix)",
			script:          "echo hello world",
			timeout:         1,
			interpreter:     "cmd",
			wantSuccessBody: "hello world\r\n",
			wantExitCode:    0,
		},
		{
			name:         "Cmd Script timeout",
			script:       "ping -n 5 127.0.0.1",
			timeout:      1,
			interpreter:  "cmd",
			wantError:    "script killed, timeout 1 sec exceeded",
			wantExitCode: 0,
		},
		{
			name: "Cmd Script succeeded (multi line)",
			script: `echo line1
echo line2`,
			timeout:         1,
			interpreter:     "cmd",
			wantSuccessBody: "line1\r\nline2\r\n",
			wantExitCode:    0,
		},
		{
			name:         "Interpreter not found",
			script:       "hostname",
			timeout:      1,
			interpreter:  "foo.exe",
			wantError:    "error starting the script: exec: \"foo.exe\": executable file not found in %PATH%",
			wantExitCode: 0,
		},
		{
			name:            "Bash on Windows",
			script:          "echo hello world",
			interpreter:     "C:\\Program Files\\Git\\bin\\bash.exe",
			wantSuccessBody: "hello world\n",
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
						Cwd:         dir,
					},
				},
			}

			//Create the actioner that implements the action interface
			var actioner actions.Actioner = runscript.Script{}

			// Execute the action by calling the mandatory function Execute()
			actionResp, err := actioner.Execute(rule, reqData)

			t.Logf("Stdout: %s\n", actionResp.SuccessBody)
			t.Logf("Stderr: %s\n", actionResp.ErrorBody)
			if tc.wantError == "" {
				assert.NoError(t, err)
				assert.Contains(t, actionResp.SuccessBody, tc.wantSuccessBody, "invalid success body")
				assert.Contains(t, actionResp.ErrorBody, tc.wantErrorBody, "invalid error body")
				assert.Equal(t, tc.wantExitCode, actionResp.Code, "exit code doesn't match")
			} else {
				assert.EqualError(t, err, tc.wantError, "Script Error")
			}
		})
	}
}
