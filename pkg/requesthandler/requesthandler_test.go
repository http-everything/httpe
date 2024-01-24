package requesthandler_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"http-everything/httpe/pkg/requesthandler"
	"http-everything/httpe/pkg/rules"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHandler(t *testing.T) {
	rule := rules.Rule{
		On: &rules.On{
			Path: "/",
		},
		Do: &rules.Do{
			AnswerContent: "foo",
		},
	}

	req, err := http.NewRequest("get", "/", nil)
	require.NoError(t, err)
	rec := httptest.NewRecorder()
	httpHandler := requesthandler.Execute(rule, nil)
	httpHandler.ServeHTTP(rec, req)

	assert.Equal(t, "foo", rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequestHandlerWithAuth(t *testing.T) {
	rule := rules.Rule{
		On: &rules.On{
			Path: "/",
		},
		Do: &rules.Do{
			AnswerContent: "foo",
		},
		With: &rules.With{
			AuthBasic: []rules.User{
				{
					Username: "john.doe",
					Password: "1234abc",
				},
			},
		},
	}

	t.Run("Access denied", func(t *testing.T) {
		// Test access denied
		req, err := http.NewRequest("get", "/", nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		httpHandler := requesthandler.Execute(rule, nil)
		httpHandler.ServeHTTP(rec, req)

		assert.Equal(t, "Unauthorised\n", rec.Body.String())
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Access granted", func(t *testing.T) {
		// Test access granted
		req, err := http.NewRequest("get", "/", nil)
		req.SetBasicAuth("john.doe", "1234abc")
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		httpHandler := requesthandler.Execute(rule, nil)
		httpHandler.ServeHTTP(rec, req)

		assert.Equal(t, "foo", rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}
