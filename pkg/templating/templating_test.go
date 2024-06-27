package templating_test

import (
	"fmt"
	"testing"

	"github.com/http-everything/httpe/pkg/templating"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/requestdata"
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
	response, err := templating.RenderActionResponse(actionResp, actionResp.SuccessBody, reqData)
	require.NoError(t, err)

	// Validate output
	want := "Hello Field Value 1"
	assert.Equal(t, want, response)
}

func TestRenderString(t *testing.T) {
	cases := []struct {
		template string
		want     string
	}{
		{
			template: "Hello {{ .Input.Form.Field1 | ToUpper}}",
			want:     "Hello FIELD VALUE 1",
		},
		{
			template: "Hello {{ .Input.Form.Field1 | ToLower}}",
			want:     "Hello field value 1",
		},
		{
			template: "Hello {{ .Input.JSON.jkey1 }}",
			want:     "Hello json value 1",
		},
		{
			template: "{{ .Input.Form.Nonexistent }}",
			want:     "",
		},
		{
			template: `{{ .Input.Form.Nonexistent|Default "John" }}`,
			want:     "John",
		},
		{
			template: `{{ .Input.JSON.nested.nkey1 }}`,
			want:     "nvalue1",
		},
		{
			template: `{{ .Input.JSON.nested.nonexistent }}`,
			want:     "<no value>",
		},
		{
			template: `{{ .Input.JSON.nonexistent }}`,
			want:     "<no value>",
		},
		{
			template: `{{ .Input.JSON.nested.nonexistent|Default "John" }}`,
			want:     "John",
		},
		{
			template: "{{ .Meta.UserAgent }}",
			want:     "golang",
		},
		{
			template: "{{ .Meta.URL }}",
			want:     "/some/path",
		},
		{
			template: "{{ .Meta.RemoteAddr }}",
			want:     "127.0.0.1",
		},
		{
			template: "{{ .Meta.Method }}",
			want:     "get",
		},
		{
			template: `{{ .Meta.Headers.Upper }}`,
			want:     "upper",
		},
		{
			template: `{{ index .Meta.Headers "X-My-Header" }}`,
			want:     "gotest",
		},
	}
	// Create mock request data
	reqData, err := requestdata.Mock()
	require.NoError(t, err)

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			// Render template
			output, err := templating.RenderString(tc.template, reqData)
			require.NoError(t, err)

			// Validate output
			assert.Equal(t, tc.want, output)
		})
	}
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
