package redirect

import (
	"net/http"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/firstof"
	"http-everything/httpe/pkg/templating"
)

type Redirect struct{}

func (r Redirect) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	var httpStatus = http.StatusFound // 302 aka Found redirects are temporary
	if rule.Action() == rules.RedirectPermanent {
		httpStatus = http.StatusMovedPermanently
	}
	location, err := templating.RenderString(
		firstof.String(rule.Do.RedirectTemporary, rule.Do.RedirectPermanent),
		reqData,
	)
	if err != nil {
		return actions.ActionResponse{}, nil
	}

	return actions.ActionResponse{
		SuccessHTTPStatus: httpStatus,
		SuccessHeaders:    map[string]string{"Location": location},
	}, nil
}
