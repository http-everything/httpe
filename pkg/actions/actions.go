package actions

import "http-everything/httpe/pkg/rules"

type ActionResponse struct {
	SuccessBody string
	ErrorBody   string
	Code        int
}

type Action interface {
	Execute(rule rules.Rule) (response ActionResponse, err error)
}
