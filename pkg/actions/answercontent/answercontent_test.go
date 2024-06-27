package answercontent_test

import (
	"fmt"
	"testing"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/answercontent"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnswerContentExecute(t *testing.T) {
	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	rule := rules.Rule{
		AnswerContent: "Agent is {{ .Meta.UserAgent }}, FormField1 is {{ .Input.Form.Field1 }}",
	}
	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answercontent.AnswerContent{}
	t.Run("no errors", func(t *testing.T) {
		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t,
			fmt.Sprintf("Agent is %s, FormField1 is %s", reqData.Meta.UserAgent, reqData.Input.Form["Field1"]),
			actionResp.SuccessBody,
		)
	})
	t.Run("with errors", func(t *testing.T) {
		rule.AnswerContent = "{{ .bad }}"
		_, err := actioner.Execute(rule, reqData)

		assert.ErrorContains(t,
			err,
			"template: simple_string:1:3: executing \"simple_string\" at <.bad>: can't evaluate "+
				"field bad in type templating.templateData",
		)
	})
}
