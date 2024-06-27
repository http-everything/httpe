package answerfile

import (
	"errors"
	"os"

	"http-everything/httpe/pkg/templating"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
)

type AnswerFile struct{}

func (n AnswerFile) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	fileContent, err := os.ReadFile(rule.AnswerFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return actions.ActionResponse{
				ErrorBody: err.Error() + "\n",
				Code:      404,
			}, nil
		}
		return actions.ActionResponse{}, err
	}

	var content string
	if rule.Args.Templating {
		content, err = templating.RenderString(string(fileContent), reqData)
		if err != nil {
			return actions.ActionResponse{}, err
		}
	} else {
		content = string(fileContent)
	}
	return actions.ActionResponse{
		SuccessBody: content,
		Code:        0,
	}, nil
}
