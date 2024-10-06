package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/firstof"
	"github.com/http-everything/httpe/pkg/share/logger"
	"github.com/http-everything/httpe/pkg/share/merge"
	"github.com/http-everything/httpe/pkg/templating"

	humanise "github.com/dustin/go-humanize" //nolint:misspell
)

const (
	DefaultOnSuccessTemplate = "{{.Action.SuccessBody }}"
	DefaultOnErrorTemplate   = "{{.Action.ErrorBody }}"
)

var DefaultHeaders map[string]string

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
	csp := map[string]string{
		"default-src": "'self'",
		"script-src":  "'self' 'nonce-2a0f584a448239d92e65e67b37264fa8' 'unsafe-eval'",
		"style-src":   "'self' 'unsafe-inline'",
		"img-src":     "'self' data:",
		"connect-src": "'self'",
		"font-src":    "'self'",
		"frame-src":   "'self'",
	}
	DefaultHeaders = make(map[string]string)
	DefaultHeaders["Content-Security-Policy"] = cspToString(csp)
	DefaultHeaders["Strict-Transport-Security"] = "max-age=63072000; includeSubDomains; preload"
	DefaultHeaders["X-Frame-Options"] = "sameorigin"
	DefaultHeaders["X-Content-Type-Options"] = "nosniff"
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
	r.w.Header().Set("WWW-Authenticate", "Basic realm='Authorization required'")
	http.Error(r.w, "Unauthorised", http.StatusUnauthorized)
}

func (r *Response) RequestEntityTooLarge(current int, limit int) {
	msg := fmt.Sprintf(
		"Request entity too large. %s sent exceeds limit of %s",
		humanise.Bytes(uint64(current)), // #nosec G115
		humanise.Bytes(uint64(limit)),   // #nosec G115
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
	response, err := templating.RenderActionResponse(actionResp, tpl, r.reqData)
	if err != nil {
		r.InternalServerError(err)
	}
	// Set all http response headers giving precedence to the headers defined by the rules over the headers set by action.
	for h, v := range merge.StringMapsI(ruleHeaders, headers, DefaultHeaders) {
		r.w.Header().Set(h, v)
	}
	r.w.WriteHeader(statusCode)
	fmt.Fprint(r.w, response)
}

func cspToString(csp map[string]string) (result string) {
	var pairs []string

	// Iterate over the map and format the key-value pairs
	for key, value := range csp {
		pairs = append(pairs, fmt.Sprintf("%s %s", key, value))
	}

	return strings.Join(pairs, "; ")
}
