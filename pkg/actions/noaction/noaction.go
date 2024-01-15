package noaction

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/rules"
)

type Noaction struct{}

func (n Noaction) Execute(_ rules.Rule) (response actions.ActionResponse, err error) {
	return actions.ActionResponse{}, nil
}
