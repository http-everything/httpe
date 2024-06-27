package middleware_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"http-everything/httpe/pkg/middleware"
	"http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func DummyRequestHandler(t *testing.T) http.Handler {
	t.Helper()
	fn := func(_ http.ResponseWriter, _ *http.Request) {}
	return http.HandlerFunc(fn)
}

func TestRequestHandlerWithAuth(t *testing.T) {
	rule := rules.Rule{
		On: &rules.On{
			Path: "/",
		},
		AnswerContent: "foo",
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
		m := middleware.New(rule, nil)
		httpHandler := m.Collection(DummyRequestHandler(t))
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
		m := middleware.New(rule, nil)
		httpHandler := m.Collection(DummyRequestHandler(t))
		httpHandler.ServeHTTP(rec, req)
		httpHandler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestRequestHandlerBodyTooLarge(t *testing.T) {
	rule := rules.Rule{
		On: &rules.On{
			Path: "/",
		},
		AnswerContent: "foo",
		With: &rules.With{
			MaxRequestBody: "4B",
		},
	}
	// Create buffer for body
	body := &bytes.Buffer{}

	// Create multipart writer
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("name", "12345")

	req, err := http.NewRequest("get", "/", body)

	require.NoError(t, err)
	rec := httptest.NewRecorder()
	m := middleware.New(rule, nil)
	httpHandler := m.Collection(DummyRequestHandler(t))
	httpHandler.ServeHTTP(rec, req)

	assert.Equal(t, "Request entity too large. 116 B sent exceeds limit of 4 B\n", rec.Body.String())
	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
}
