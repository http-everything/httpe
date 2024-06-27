package renderbuttons

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/templating"
)

//go:embed buttons.tpl.html
var buttonsTpl string

type RenderButtons struct{}

type buttons struct {
	Title   string
	Buttons []rules.Button
}

func (r RenderButtons) Execute(rule rules.Rule, _ requestdata.Data) (response actions.ActionResponse, err error) {
	if rule.Args.Template != "" {
		// Overwrite the embedded template with a file from the file system
		t, err := os.ReadFile(rule.Args.Template)
		if err != nil {
			return actions.ActionResponse{}, fmt.Errorf("error reading template: %w", err)
		}
		buttonsTpl = string(t)
	}
	te, err := template.New("buttons").Funcs(templating.TplFuncs).Parse(buttonsTpl)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	var html bytes.Buffer
	err = te.Execute(&html, buttons{
		Title:   rule.Name,
		Buttons: rule.RenderButtons,
	})
	if err != nil {
		return actions.ActionResponse{}, err
	}

	return actions.ActionResponse{
		SuccessBody: html.String(),
		Code:        0,
	}, nil
}
