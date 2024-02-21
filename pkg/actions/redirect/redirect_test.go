package redirect_test

import (
	"fmt"
	"net/http"
	"testing"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/redirect"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectExecute(t *testing.T) {
	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	var actioner actions.Actioner = redirect.Redirect{}
	t.Run("Redir Perm", func(t *testing.T) {
		rule := rules.Rule{
			Do: &rules.Do{
				RedirectPermanent: "{{ .Input.URLPlaceholders.redir }}",
			},
		}
		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf(reqData.Input.URLPlaceholders["redir"]), actionResp.SuccessHeaders["Location"])
		assert.Equal(t, http.StatusMovedPermanently, actionResp.SuccessHTTPStatus)
	})
	t.Run("Redir Temp", func(t *testing.T) {
		rule := rules.Rule{
			Do: &rules.Do{
				RedirectTemporary: "{{ .Input.URLPlaceholders.redir }}",
			},
		}
		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf(reqData.Input.URLPlaceholders["redir"]), actionResp.SuccessHeaders["Location"])
		assert.Equal(t, http.StatusFound, actionResp.SuccessHTTPStatus)
	})
}
