package actions

import (
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
)

type ActionResponse struct {
	SuccessBody       string            `json:"success_body"`
	ErrorBody         string            `json:"error_body"`
	Code              int               `json:"code"`
	SuccessHTTPStatus int               `json:"-"`
	ErrorHTTPStatus   int               `json:"-"`
	SuccessHeaders    map[string]string `json:"-"`
	ErrorHeaders      map[string]string `json:"-"`
}

type Actioner interface {
	Execute(rule rules.Rule, reqData requestdata.Data) (response ActionResponse, err error)
}
