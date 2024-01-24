package actions

import (
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
)

type ActionResponse struct {
	SuccessBody string
	ErrorBody   string
	Code        int
}

type Actioner interface {
	Execute(rule rules.Rule, reqData requestdata.Data) (response ActionResponse, err error)
}
