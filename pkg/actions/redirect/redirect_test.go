package redirect_test

import (
	"net/http"
	"testing"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/redirect"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectExecute(t *testing.T) {
	reqData, err := requestdata.Mock()
	require.NoError(t, err)
	var actioner actions.Actioner = redirect.Redirect{}
	t.Run("Redir Perm", func(t *testing.T) {
		rule := rules.Rule{
			RedirectPermanent: "{{ .Input.URLPlaceholders.redir }}",
		}
		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t, reqData.Input.URLPlaceholders["redir"], actionResp.SuccessHeaders["Location"])
		assert.Equal(t, http.StatusMovedPermanently, actionResp.SuccessHTTPStatus)
	})
	t.Run("Redir Temp", func(t *testing.T) {
		rule := rules.Rule{
			RedirectTemporary: "{{ .Input.URLPlaceholders.redir }}",
		}
		actionResp, err := actioner.Execute(rule, reqData)
		require.NoError(t, err)

		assert.Equal(t, reqData.Input.URLPlaceholders["redir"], actionResp.SuccessHeaders["Location"])
		assert.Equal(t, http.StatusFound, actionResp.SuccessHTTPStatus)
	})
}
