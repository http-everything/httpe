package answerfile

import (
	"errors"
	"os"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/templating"
)

type AnswerFile struct{}

func (n AnswerFile) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	fileContent, err := os.ReadFile(rule.Do.AnswerFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return actions.ActionResponse{
				ErrorBody: err.Error() + "\n",
				Code:      404,
			}, nil
		}
		return actions.ActionResponse{}, err
	}

	content, err := templating.RenderString(string(fileContent), reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	return actions.ActionResponse{
		SuccessBody: content,
		Code:        0,
	}, nil
}
