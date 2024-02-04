package templating

import (
	"bytes"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"io"
	"reflect"
	"strings"
	"text/template"
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
	"Default": func(arg interface{}, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				return arg
			}
		case reflect.Bool:
			if !v.Bool() {
				return arg
			}
		default:
			return value
		}

		return value
	},
}

func RenderActionResponse(actionResp actions.ActionResponse, tpl string, reqData requestdata.Data, wr io.Writer) (err error) {
	te, err := template.New("response").Funcs(TplFuncs).Parse(tpl)
	if err != nil {
		return err
	}
	tplData := templateData{
		Action: actionResp,
		Meta:   reqData.Meta,
		Input:  reqData.Input,
	}
	err = te.Execute(wr, tplData)
	if err != nil {
		return err
	}
	return nil
}

func RenderString(input string, reqData requestdata.Data) (output string, err error) {
	te, err := template.New("input").Funcs(TplFuncs).Parse(input)
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
