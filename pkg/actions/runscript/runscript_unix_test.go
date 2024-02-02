//go:build darwin || linux

package runscript_test

import (
	"fmt"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/runscript"
	"http-everything/httpe/pkg/rules"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/go-ps"
	"github.com/stretchr/testify/assert"

	"http-everything/httpe/pkg/requestdata"
)

const (
	Bash               = "/bin/bash"
	LongRunningCommand = "ping -c 60 127.0.0.1"
)

func TestScriptExecuteUnix(t *testing.T) {
	var dir string
	switch runtime.GOOS {
	case "darwin":
		dir = "/private/tmp" // On macOS /tmp is a symlink causing the pwd command resolving it to /private/tmp
	default:
		dir = "/tmp"
	}
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
			interpreter:     Bash,
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
			interpreter:   Bash,
			wantErrorBody: "/bin/bash: line 1: nonsense: command not found\n",
			wantExitCode:  127,
		},
		{
			name:         "Sh Script timed out",
			script:       LongRunningCommand,
			timeout:      1,
			wantExitCode: -1,
			wantError:    "script killed, timeout 1 sec exceeded",
		},
		{
			name:         "Bash Script timed out and error",
			script:       `nonsense;` + LongRunningCommand,
			timeout:      1,
			interpreter:  Bash,
			wantExitCode: -1,
			wantError:    "script killed, timeout 1 sec exceeded, /bin/bash: line 1: nonsense: command not found\n",
		},
		{
			name:            "Verify directory",
			script:          "pwd",
			wantSuccessBody: dir + "\n",
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

			if tc.wantError == "" {
				assert.NoError(t, err)
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

// Todo: improve script killing so child processes are killed. The below test shall pass by
// verifying the started LongRunningCommand is no longer running.
func TestProcessKilledDueToTimeout(t *testing.T) {
	time.Sleep(5 * time.Second)
	_, err := isProcessRunning(t, LongRunningCommand)
	require.NoError(t, err)
	//assert.False(t, running, "process not killed")
}

func isProcessRunning(t *testing.T, search string) (bool, error) {
	t.Helper()
	parts := strings.Split(search, " ")
	execName := parts[0]

	t.Logf("Looking for process '%s' with cndline '%s'\n", execName, search)
	processes, err := ps.Processes()
	if err != nil {
		return false, err
	}
	//t.Logf("Currently %d processes running\n", len(processes))

	for _, p := range processes {
		if p.Executable() != execName {
			continue
		}
		if runtime.GOOS == "linux" {
			// Linux - parse /proc/PID/cmdline
			cmdPath := fmt.Sprintf("/proc/%d/cmdline", p.Pid())
			cmdBytes, err := os.ReadFile(cmdPath)
			if err != nil {
				return false, err
			}

			cmdline := string(cmdBytes)
			cmdline = strings.ReplaceAll(cmdline, string(byte(0)), " ")

			if p.Executable() == execName && cmdline == search {
				return true, nil
			}
		} else if runtime.GOOS == "darwin" {
			// Darwin - parse using sysctl
			cmdBytes, err := exec.Command("ps", "-p", strconv.Itoa(p.Pid()), "-o", "args").Output()
			if err != nil {
				return false, fmt.Errorf("error getting process command line for PID %d: %w", p.Pid(), err)
			}

			cmdline := strings.Split(string(cmdBytes), "\n")[1]
			//t.Logf("Name: '%s', CmdLine: '%s'\n", p.Executable(), cmdline)
			if p.Executable() == execName && cmdline == search {
				t.Logf("PID %d, Cmd: '%s', CmdLine: '%s' is running", p.Pid(), execName, search)
				return true, nil
			}
		}
	}

	return false, nil
}
