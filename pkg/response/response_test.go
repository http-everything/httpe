package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/response"
	"http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalServerError(t *testing.T) {
	// Create a ResponseRecorder (which satisfies http.ResponseWriter)
	rr := httptest.NewRecorder()
	ruleResp := rules.Respond{}

	resp := response.New(rr, ruleResp, nil)
	resp.InternalServerError(errors.New("error"))

	// Validate status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "error\n", rr.Body.String())
}

func TestInternalServerErrorf(t *testing.T) {
	// Create a ResponseRecorder (which satisfies http.ResponseWriter)
	rr := httptest.NewRecorder()
	ruleResp := rules.Respond{}

	resp := response.New(rr, ruleResp, nil)
	resp.InternalServerErrorf("error %s", "test")

	// Validate status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "error test\n", rr.Body.String())
}

func TestActionResponse(t *testing.T) {
	tplAllFields, err := os.ReadFile("../../testdata/templates/all-fields.tpl")
	require.NoError(t, err)
	var jd interface{}
	err = json.Unmarshal([]byte(`{"jsonkey1": "jsonvalue1"}`), &jd)
	require.NoError(t, err)
	cases := []struct {
		name       string
		acRes      actions.ActionResponse
		ruRes      rules.Respond
		reqData    *requestdata.Data
		want       []string
		wantStatus int
	}{
		{
			name: "should return the success body using the default template",
			acRes: actions.ActionResponse{
				SuccessBody: "all good",
				Code:        0,
			},
			ruRes:      rules.Respond{},
			reqData:    nil,
			want:       []string{"all good"},
			wantStatus: http.StatusOK,
		},
		{
			name: "should return the error body using the default template",
			acRes: actions.ActionResponse{
				ErrorBody: "nothing is good",
				Code:      1,
			},
			ruRes:      rules.Respond{},
			reqData:    nil,
			want:       []string{"nothing is good"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "should return body from custom template on success",
			acRes: actions.ActionResponse{
				SuccessBody: "good",
				ErrorBody:   "bad",
				Code:        0,
			},
			ruRes: rules.Respond{
				OnSuccess: rules.OnSuccess{
					Body: "test {{.Action.SuccessBody }} {{.Action.ErrorBody }} test",
				},
			},
			reqData:    nil,
			want:       []string{"test good bad test"},
			wantStatus: http.StatusOK,
		},
		{
			name: "should return body from custom template on error",
			acRes: actions.ActionResponse{
				SuccessBody: "good",
				ErrorBody:   "bad",
				Code:        99,
			},
			ruRes: rules.Respond{
				OnError: rules.OnError{
					Body:       "test {{.Action.SuccessBody }} {{.Action.ErrorBody }} {{.Action.Code }} test",
					HTTPStatus: http.StatusConflict,
				},
			},
			reqData:    nil,
			want:       []string{"test good bad 99 test"},
			wantStatus: http.StatusConflict,
		},
		{
			name: "should return body from custom template on success with request data",
			acRes: actions.ActionResponse{
				SuccessBody: "good",
				ErrorBody:   "bad",
				Code:        0,
			},
			ruRes: rules.Respond{
				OnSuccess: rules.OnSuccess{
					Body: string(tplAllFields),
				},
			},
			reqData: &requestdata.Data{
				Meta: requestdata.MetaData{
					UserAgent: "Test Agent",
					Headers:   map[string]string{"X-Test-Header": "foo"},
					Method:    "GET",
				},
				Input: requestdata.Input{
					Form: map[string]string{
						"formkey1": "formvalue1",
					},
					Params: map[string]string{
						"paramkey1": "paramvalue1",
					},
					Uploads: []requestdata.Upload{
						{
							FieldName: "myfile",
							FileName:  "myfile.txt",
							Size:      100,
							Type:      "text/plain",
							Stored:    "/tmp/myfile.txt",
						},
					},
					JSON: jd,
				},
			},
			want: []string{
				"formkey1: formvalue1",
				"paramkey1: paramvalue1",
				"Field Name: myfile",
				"File Name: myfile.txt",
				"Stored: /tmp/myfile.txt",
				"Size: 100",
				"User Agent: Test Agent",
				"X-Test-Header: foo",
				"Method: GET",
				"jsonkey1: jsonvalue1",
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a ResponseRecorder (which satisfies http.ResponseWriter)
			rr := httptest.NewRecorder()
			resp := response.New(rr, tc.ruRes, nil)
			if tc.reqData != nil {
				resp.AddRequestData(*tc.reqData)
			}
			resp.ActionResponse(tc.acRes)

			// Validate status code
			assert.Equal(t, tc.wantStatus, rr.Code)
			for _, w := range tc.want {
				assert.Contains(t, rr.Body.String(), w)
			}
		})
	}
}
