package answerfile_test

import (
	"fmt"
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/answerfile"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnswerFileExecute(t *testing.T) {
	testFile := t.TempDir() + "/test.txt"
	err := os.WriteFile(testFile, []byte("Agent is {{ .Meta.UserAgent }}, FormField1 is {{ .Input.Form.Field1 }}"), 0400)
	require.NoError(t, err)

	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	rule := rules.Rule{
		Do: &rules.Do{
			AnswerFile: testFile,
		},
	}

	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answerfile.AnswerFile{}
	actionResp, err := actioner.Execute(rule, reqData)
	require.NoError(t, err)

	assert.Equal(t,
		fmt.Sprintf("Agent is %s, FormField1 is %s", reqData.Meta.UserAgent, reqData.Input.Form["Field1"]),
		actionResp.SuccessBody,
	)
}

func TestAnswerFileExecuteFileNotFound(t *testing.T) {
	testFile := t.TempDir() + "/test.txt"
	rule := rules.Rule{
		Do: &rules.Do{
			AnswerFile: testFile,
		},
	}
	//Create the actioner that implements the action interface
	var actioner actions.Actioner = answerfile.AnswerFile{}
	actionResp, err := actioner.Execute(rule, requestdata.Data{})
	require.NoError(t, err)

	assert.Equal(t, "", actionResp.SuccessBody)
	assert.Contains(t, actionResp.ErrorBody, "test.txt: no such file or directory")
}
