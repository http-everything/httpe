package answerfile_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/answerfile"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnswerFileExecute(t *testing.T) {
	testFile := t.TempDir() + "/test.txt"

	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	rule := rules.Rule{
		AnswerFile: testFile,
	}

	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answerfile.AnswerFile{}
	t.Run("no errors, templating off", func(t *testing.T) {
		err := os.WriteFile(testFile, []byte("This is a {{ .String }}"), 0600)
		require.NoError(t, err)

		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t, "This is a {{ .String }}", actionResp.SuccessBody)
	})

	t.Run("no errors, templating on", func(t *testing.T) {
		rule.Args.Templating = true
		err := os.WriteFile(testFile, []byte("Agent is {{ .Meta.UserAgent }}, FormField1 is {{ .Input.Form.Field1 }}"), 0600)
		require.NoError(t, err)

		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t,
			fmt.Sprintf("Agent is %s, FormField1 is %s", reqData.Meta.UserAgent, reqData.Input.Form["Field1"]),
			actionResp.SuccessBody,
		)
	})

	t.Run("with errors", func(t *testing.T) {
		err := os.WriteFile(testFile, []byte("{{ .bad }}"), 0600)
		require.NoError(t, err)

		_, err = actioner.Execute(rule, reqData)

		assert.ErrorContains(t, err, "template: simple_string:1:3: executing \"simple_string\" at <.bad>: "+
			"can't evaluate field bad in type templating.templateData")
	})
}

func TestAnswerFileNotFound(t *testing.T) {
	testFile := t.TempDir() + "/test.txt"
	rule := rules.Rule{
		AnswerFile: testFile,
	}
	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answerfile.AnswerFile{}
	actionResp, err := actioner.Execute(rule, requestdata.Data{})
	require.NoError(t, err)

	assert.Equal(t, "", actionResp.SuccessBody)
	if runtime.GOOS == "windows" {
		assert.Contains(t, actionResp.ErrorBody, "test.txt: The system cannot find the file specified")
	} else {
		assert.Contains(t, actionResp.ErrorBody, "test.txt: no such file or directory")
	}
}
