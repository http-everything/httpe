package sendemail

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"net/http"
)

type Email struct{}

func (e Email) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	return actions.ActionResponse{
		SuccessBody:       "not implemented yet",
		SuccessHTTPStatus: http.StatusNotImplemented,
	}, nil
}
