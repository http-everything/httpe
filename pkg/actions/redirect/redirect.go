package redirect

import (
	"net/http"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/firstof"
	"github.com/http-everything/httpe/pkg/templating"
)

type Redirect struct{}

func (r Redirect) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	var httpStatus = http.StatusFound // 302 aka Found redirects are temporary
	if rule.Action() == rules.RedirectPermanent {
		httpStatus = http.StatusMovedPermanently
	}
	location, err := templating.RenderString(
		firstof.String(rule.RedirectTemporary, rule.RedirectPermanent),
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
