package response

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/firstof"
	"http-everything/httpe/pkg/share/logger"
	"http-everything/httpe/pkg/share/merge"
	"http-everything/httpe/pkg/templating"
	"net/http"
)

const (
	DefaultOnSuccessTemplate = "{{.Action.SuccessBody }}"
	DefaultOnErrorTemplate   = "{{.Action.ErrorBody }}"
)

type Response struct {
	reqData  requestdata.Data
	logger   *logger.Logger
	w        http.ResponseWriter
	ruleResp rules.Respond
}

func New(w http.ResponseWriter, ruleResp rules.Respond, logger *logger.Logger) *Response {
	resp := &Response{
		reqData:  requestdata.Data{},
		w:        w,
		logger:   logger,
		ruleResp: ruleResp,
	}
	return resp
}

func (r *Response) AddRequestData(reqData requestdata.Data) {
	r.reqData = reqData
}

func (r *Response) InternalServerError(err error) {
	if r.logger != nil {
		r.logger.Errorf("Internal Server Error: %v", err)
	}
	http.Error(r.w, err.Error(), http.StatusInternalServerError)
}

func (r *Response) Unauthorised() {
	http.Error(r.w, "Unauthorised", http.StatusUnauthorized)
}

func (r *Response) RequestEntityTooLarge(current int, limit int) {
	msg := fmt.Sprintf(
		"Request entity too large. %s sent exceeds limit of %s",
		humanize.Bytes(uint64(current)),
		humanize.Bytes(uint64(limit)),
	)
	http.Error(r.w, msg, http.StatusRequestEntityTooLarge)
}

func (r *Response) InternalServerErrorf(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Errorf("Internal Server Error: " + msg)
	}
	http.Error(r.w, msg, http.StatusInternalServerError)
}

func (r *Response) ActionResponse(actionResp actions.ActionResponse) {
	var tpl string
	var statusCode int
	var headers map[string]string
	var actionSucceeded = false
	// Set a template giving what's defined in the rule precedence over the default defined by the action response
	if actionResp.Code != 0 {
		// Handle a failed action
		// Set a template that later will render the error response
		tpl = firstof.String(r.ruleResp.OnError.Body, DefaultOnErrorTemplate)
		// Set the HTTP Status code giving what's defined in the rule precedence over the default defined by the action response
		statusCode = firstof.Int(
			r.ruleResp.OnError.HTTPStatus,
			actionResp.ErrorHTTPStatus,
			http.StatusBadRequest,
		)
		headers = actionResp.ErrorHeaders
	} else {
		// Handle the succeeded action
		actionSucceeded = true
		// Set a template that later will render the success response
		tpl = firstof.String(r.ruleResp.OnSuccess.Body, DefaultOnSuccessTemplate)
		// Set the HTTP Status code giving what's defined in the rule precedence over the default defined by the action response
		statusCode = firstof.Int(
			r.ruleResp.OnSuccess.HTTPStatus,
			actionResp.SuccessHTTPStatus,
			http.StatusOK,
		)
		headers = actionResp.SuccessHeaders
	}
	// Render the http header defined by the rule.
	ruleHeaders, err := templating.RenderStringMap(r.ruleResp.Headers(actionSucceeded), r.reqData)
	if err != nil {
		r.InternalServerError(err)
		return
	}
	// Set all http response headers giving precedence to the headers defined by the rules over the headers set by action.
	for h, v := range merge.StringMapsI(ruleHeaders, headers) {
		r.w.Header().Set(h, v)
	}
	r.w.WriteHeader(statusCode)
	err = templating.RenderActionResponse(actionResp, tpl, r.reqData, r.w)
	if err != nil {
		r.InternalServerError(err)
	}
}
