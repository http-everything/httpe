package templating_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"http-everything/httpe/pkg/templating"
	"testing"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
)

func TestRenderActionResponse(t *testing.T) {
	// Create mock request data
	reqData, err := requestdata.Mock()
	require.NoError(t, err)

	// Create mock action response
	actionResp := actions.ActionResponse{
		SuccessBody: "Hello {{.Input.Form.Field1}}",
	}

	// Render template
	var buf bytes.Buffer
	err = templating.RenderActionResponse(actionResp, actionResp.SuccessBody, reqData, &buf)
	require.NoError(t, err)

	// Validate output
	want := "Hello Field Value 1"
	assert.Equal(t, want, buf.String())
}

func TestRenderString(t *testing.T) {
	// Create mock request data
	reqData, err := requestdata.Mock()
	require.NoError(t, err)

	// Render template
	inputTpl := "Hello {{.Input.Form.Field1 | ToUpper}}"
	output, err := templating.RenderString(inputTpl, reqData)
	require.NoError(t, err)

	// Validate output
	want := "Hello FIELD VALUE 1"
	assert.Equal(t, want, output)
}

func TestRenderStringMap(t *testing.T) {
	// Create mock request data
	reqData, err := requestdata.Mock()
	require.NoError(t, err)

	// Render template
	inputTpl := map[string]string{
		"key1": "Hello {{.Input.Form.Field1 | ToUpper}}",
		"key2": "Hello {{.Input.Form.Field2 | ToUpper}}",
	}
	output, err := templating.RenderStringMap(inputTpl, reqData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Validate output
	want := map[string]string{
		"key1": "Hello FIELD VALUE 1",
		"key2": "Hello FIELD VALUE 2",
	}
	assert.Equal(t, want, output)
}
