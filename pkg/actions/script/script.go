package script

import (
	"bytes"
	"errors"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/set"
	"io"
	"os/exec"
)

const (
	DefaultInterpreter = "sh"
)

type Script struct{}

func (s Script) Execute(rule rules.Rule) (response actions.ActionResponse, err error) {
	var exitCode = 0
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(set.String(rule.Do.Args.Interpreter, DefaultInterpreter))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	stdin, _ := cmd.StdinPipe()
	_, err = io.WriteString(stdin, rule.Do.Script)
	if err != nil {
		return actions.ActionResponse{}, err
	}
	err = stdin.Close()
	if err != nil {
		return actions.ActionResponse{}, err
	}
	err = cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		}
	}
	return actions.ActionResponse{
		SuccessBody: stdout.String(),
		ErrorBody:   stderr.String(),
		Code:        exitCode,
	}, nil
}
