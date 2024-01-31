package runscript

import (
	"context"
	"errors"
	"fmt"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/buffer"
	"http-everything/httpe/pkg/share/firstof"
	"http-everything/httpe/pkg/share/remove"
	"http-everything/httpe/pkg/templating"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v4"
)

const (
	DefaultTimeoutSecs = 30
)

type Script struct{}

func (s Script) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	var exitCode = -1
	timeoutSec := firstof.Int(rule.Do.Args.Timeout, DefaultTimeoutSecs)
	script, _ := templating.RenderString(rule.Do.RunScript, reqData)
	interpreter := firstof.String(rule.Do.Args.Interpreter, defaultInterpreter())

	// Change to the working directory if specified by the rule.
	err = os.Chdir(firstof.String(rule.Do.Args.Cwd, os.TempDir()))
	if err != nil {
		return actions.ActionResponse{}, fmt.Errorf("error changing to directory '%s': %w", rule.Do.Args.Cwd, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	args, cleanupFn, err := defaultInterpreterArgs(interpreter, script)
	if err != nil {
		return actions.ActionResponse{}, fmt.Errorf("error preparing interpreter: %w", err)
	}
	defer cleanupFn()

	cmd := exec.CommandContext(ctx, interpreter, args...)

	// Create a stdin pipe to the script
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return actions.ActionResponse{}, fmt.Errorf("error creating stdin pipe: %w", err)
	}

	var stdinBu buffer.Buffer
	defer stdinBu.CollectingDone()
	cmd.Stdout = &stdinBu

	var stderrBu buffer.Buffer
	defer stderrBu.CollectingDone()
	cmd.Stderr = &stderrBu

	if err = cmd.Start(); err != nil { // Use start, not run
		return actions.ActionResponse{}, fmt.Errorf("error starting the script: %w", err)
	}

	// Write the script to stdin of the interpreter
	_, err = io.WriteString(stdin, script)
	if err != nil {
		return actions.ActionResponse{}, fmt.Errorf("error writing to stdin pipe: %w", err)
	}

	err = stdin.Close()
	if err != nil {
		return actions.ActionResponse{}, fmt.Errorf("error closing stdin pipe: %w", err)
	}

	killCh := make(chan error, 10)
	doneCh := make(chan error, 10)
	exitCodeCh := make(chan int, 10)

	// Start goroutine to monitor for context timeout and kill the process.
	// goroutine waits on either the context or the command to finish
	process, processDone := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			err := cmd.Process.Kill()
			if err != nil {
				killCh <- fmt.Errorf("killing of the script failed, script killed because of ctx cancel: %w", err)
			} else {
				killCh <- fmt.Errorf("script killed")
			}

		case <-process.Done():
			exitCodeCh <- cmd.ProcessState.ExitCode()
		}
		close(exitCodeCh)
		close(killCh)
	}()

	// Start go routine to wait for completion and signal on the doneCh channel
	go func() {
		err := cmd.Wait()
		if err != nil {
			doneCh <- fmt.Errorf("process error: %w", err)
		}
		close(doneCh)
		processDone()
	}()

	var errs []string

	select {
	case <-ctx.Done():
		select {
		case <-process.Done():
			for err := range doneCh {
				errs = append(errs, err.Error())
			}
		case <-time.After(time.Millisecond * 200):
			for err := range killCh {
				errs = append(errs, err.Error())
			}
			errs = append(errs, fmt.Errorf("timeout %d sec exceeded", timeoutSec).Error())
		}
	case <-process.Done():
		for err := range doneCh {
			errs = append(errs, err.Error())
		}
		for c := range exitCodeCh {
			exitCode = c
		}
	}

	if stderrBu.String() != "" {
		errs = append(errs, stderrBu.String())
	}
	if exitCode < 0 {
		return actions.ActionResponse{}, errors.New(strings.Join(errs, ", "))
	}

	return actions.ActionResponse{
		SuccessBody: stdinBu.String(),
		ErrorBody:   stderrBu.String(),
		Code:        exitCode,
	}, nil
}

func defaultInterpreter() string {
	switch runtime.GOOS {
	case "windows":
		return "powershell"
	default:
		return "sh"
	}
}

func defaultInterpreterArgs(interpreter string, script string) (args []string, cleanupFn func(), err error) {
	interpreter = strings.ToLower(remove.FileExtension(interpreter, "exe"))
	noCleanup := func() {
		// Empty function when no clean up is required
	}
	switch interpreter {
	case "powershell", "pwsh":
		return []string{"-NoProfile", "-NonInteractive", "-Command", "-"}, noCleanup, nil
	case "cmd":
		//storing the script in a file and appending the script file to the command line appears to be the only
		//solution for cmd.exe because the /Q/C switches cause cmd.exe to ignore the stdin pipe.
		bat, err := scriptToBat(script)
		if err != nil {
			return []string{}, noCleanup, err
		}
		fn := func() {
			os.Remove(bat)
		}
		return []string{"/Q/C", bat}, fn, nil
	default:
		return []string{}, noCleanup, nil
	}
}

func scriptToBat(script string) (filepath string, err error) {
	filepath = os.TempDir() + "/httpe_" + shortuuid.New() + ".bat"
	err = os.WriteFile(filepath, []byte(script), 0600)
	if err != nil {
		return "", fmt.Errorf("error saving script to temporary file %s: %w", filepath, err)
	}

	return filepath, nil
}
