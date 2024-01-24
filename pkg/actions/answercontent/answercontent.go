package answercontent

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/templating"
)

type AnswerContent struct{}

func (n AnswerContent) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	content, err := templating.RenderActionInput(rule.Do.AnswerContent, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	return actions.ActionResponse{
		SuccessBody: content,
		Code:        0,
	}, nil
}
