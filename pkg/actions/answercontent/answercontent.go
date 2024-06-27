package answercontent

import (
	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/templating"
)

type AnswerContent struct{}

func (n AnswerContent) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	content, err := templating.RenderString(rule.AnswerContent, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	return actions.ActionResponse{
		SuccessBody: content,
		Code:        0,
	}, nil
}
