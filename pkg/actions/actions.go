package actions

import (
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
)

type ActionResponse struct {
	SuccessBody       string
	ErrorBody         string
	Code              int
	SuccessHTTPStatus int
	ErrorHTTPStatus   int
	SuccessHeaders    map[string]string
	ErrorHeaders      map[string]string
}

type Actioner interface {
	Execute(rule rules.Rule, reqData requestdata.Data) (response ActionResponse, err error)
}
