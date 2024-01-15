package response

import (
	"fmt"
	"html/template"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/request"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"
	"http-everything/httpe/pkg/share/set"
	"net/http"
)

const (
	DefaultOnSuccessTemplate = "{{.Action.SuccessBody }}"
	DefaultOnErrorTemplate   = "{{.Action.ErrorBody }}"
)

type Response struct {
	meta     request.MetaData
	logger   *logger.Logger
	w        http.ResponseWriter
	ruleResp rules.Respond
}

type templateData struct {
	Action actions.ActionResponse
	Meta   request.MetaData
}

func New(w http.ResponseWriter, ruleResp rules.Respond, meta request.MetaData, logger *logger.Logger) *Response {
	resp := &Response{
		meta:     meta,
		w:        w,
		logger:   logger,
		ruleResp: ruleResp,
	}
	return resp
}

func (r *Response) InternalServerError(err error) {
	if r.logger != nil {
		r.logger.Errorf("Internal Server Error: %v", err)
	}
	http.Error(r.w, err.Error(), http.StatusInternalServerError)
}

func (r *Response) InternalServerErrorf(msg string, args ...interface{}) {
	msg = fmt.Sprintf("Internal Server Error: "+msg, args...)
	if r.logger != nil {
		r.logger.Errorf(msg)
	}
	http.Error(r.w, msg, http.StatusInternalServerError)
}

func (r *Response) ActionResponse(actionResp actions.ActionResponse) {
	var tpl string
	var statusCode int
	// Set a template giving what's defined in the rule precedence over the default defined by the action response
	if actionResp.Code != 0 {
		// Handle a failed action
		tpl = set.String(r.ruleResp.OnError.Body, DefaultOnErrorTemplate)
		// Set the HTTP Status code giving what's defined in the rule precedence over the default defined by the action response
		statusCode = set.Int(r.ruleResp.OnError.HTTPStatus, 400)
	} else {
		// Handle the succeeded action
		tpl = set.String(r.ruleResp.OnSuccess.Body, DefaultOnSuccessTemplate)
		// Set the HTTP Status code giving what's defined in the rule precedence over the default defined by the action response
		statusCode = set.Int(r.ruleResp.OnSuccess.HTTPStatus, 200)
	}
	te, err := template.New("response").Funcs(tplFunc).Parse(tpl)
	if err != nil {
		r.InternalServerError(err)
		return
	}

	tplData := templateData{
		Action: actionResp,
		Meta:   r.meta,
	}
	r.w.WriteHeader(statusCode)
	err = te.Execute(r.w, tplData)
	if err != nil {
		r.InternalServerError(err)
	}
}
