package templating_test

import (
	"bytes"
	"http-everything/httpe/pkg/templating"
	"testing"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
)

func TestRenderActionResponse(t *testing.T) {
	// Create mock request data
	reqData := requestdata.Data{
		Meta: requestdata.MetaData{},
		Input: requestdata.Input{
			Form: map[string]string{"name": "John"},
		},
	}

	// Create mock action response
	actionResp := actions.ActionResponse{
		SuccessBody: "Hello {{.Input.Form.name}}",
	}

	// Render template
	var buf bytes.Buffer
	err := templating.RenderActionResponse(actionResp, actionResp.SuccessBody, reqData, &buf)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate output
	want := "Hello John"
	if buf.String() != want {
		t.Errorf("Expected %q, got %q", want, buf.String())
	}
}

func TestRenderActionInput(t *testing.T) {
	// Create mock request data
	reqData := requestdata.Data{
		Meta: requestdata.MetaData{},
		Input: requestdata.Input{
			Form: map[string]string{"name": "John"},
		},
	}

	// Render template
	inputTpl := "Hello {{.Input.Form.name | ToUpper}}"
	output, err := templating.RenderActionInput(inputTpl, reqData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate output
	want := "Hello JOHN"
	if output != want {
		t.Errorf("Expected %q, got %q", want, output)
	}
}
