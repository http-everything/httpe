package sendemail

import (
	"net/http"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
)

type Email struct{}

func (e Email) Execute(_ rules.Rule, _ requestdata.Data) (response actions.ActionResponse, err error) {
	return actions.ActionResponse{
		SuccessBody:       "not implemented yet",
		SuccessHTTPStatus: http.StatusNotImplemented,
	}, nil
}
