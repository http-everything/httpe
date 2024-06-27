package templating

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/requestdata"
)

type templateData struct {
	Action actions.ActionResponse
	Meta   requestdata.MetaData
	Input  requestdata.Input
}

// recovery will silently swallow all unexpected panics.
func recovery() {
	_ = recover()
}

var TplFuncs = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
	"Default": func(defVal interface{}, curVal interface{}) interface{} {
		defer recovery()

		c := reflect.ValueOf(curVal)
		if !c.IsValid() {
			return defVal
		}
		if c.IsZero() {
			return defVal
		}
		switch c.Kind() {

		case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
			if c.Len() == 0 {
				return defVal
			}
		case reflect.Bool:
			if !c.Bool() {
				return defVal
			}
		default:
			return curVal
		}
		return curVal
	},
}

func RenderActionResponse(actionResp actions.ActionResponse, tpl string, reqData requestdata.Data) (response string, err error) {
	te, err := newTpl(tpl, "action_response")
	if err != nil {
		return "", err
	}
	tplData := templateData{
		Action: actionResp,
		Meta:   reqData.Meta,
		Input:  reqData.Input,
	}
	var bu bytes.Buffer
	err = te.Execute(&bu, tplData)
	if err != nil {
		return "", err
	}
	return bu.String(), nil
}

func RenderString(input string, reqData requestdata.Data) (output string, err error) {
	te, err := newTpl(input, "simple_string")
	if err != nil {
		return "", err
	}
	tplData := templateData{
		Meta:  reqData.Meta,
		Input: reqData.Input,
	}
	var bu bytes.Buffer
	err = te.Execute(&bu, tplData)
	if err != nil {
		return "", err
	}
	return bu.String(), nil
}

func RenderStringMap(input map[string]string, reqData requestdata.Data) (output map[string]string, err error) {
	output = make(map[string]string)
	for k, v := range input {
		if output[k], err = RenderString(v, reqData); err != nil {
			return output, err
		}
	}
	return output, nil
}

func newTpl(input string, name string) (*template.Template, error) {
	return template.New(name).Funcs(TplFuncs).Option("missingkey=zero").Parse(input)
}
