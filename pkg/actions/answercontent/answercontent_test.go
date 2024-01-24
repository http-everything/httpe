package answercontent_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/answercontent"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"testing"
)

func TestAnswerContent_Execute(t *testing.T) {
	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	rule := rules.Rule{
		Do: &rules.Do{
			AnswerContent: "Agent is {{ .Meta.UserAgent }}, FormField1 is {{ .Input.Form.Field1 }}",
		},
	}
	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answercontent.AnswerContent{}
	actionResp, err := actioner.Execute(rule, reqData)
	require.NoError(t, err)

	assert.Equal(t,
		fmt.Sprintf("Agent is %s, FormField1 is %s", reqData.Meta.UserAgent, reqData.Input.Form["Field1"]),
		actionResp.SuccessBody,
	)
}
