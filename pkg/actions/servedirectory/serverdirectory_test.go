package servedirectory_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"http-everything/httpe/pkg/actions/servedirectory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(dir+"/test.html", []byte("hello world"), 0400)
	require.NoError(t, err)
	t.Run("listing", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/static", nil)
		handler := servedirectory.Handle("/static", dir)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		assert.Contains(t, w.Body.String(), "<a href=\"test.html\">test.html</a>")
	})

	t.Run("get file", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/test.html", nil)
		handler := servedirectory.Handle("", dir)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		assert.Equal(t, "hello world", w.Body.String())
	})

	t.Run("get index.html", func(t *testing.T) {
		err := os.WriteFile(dir+"/index.html", []byte("start here"), 0400)
		require.NoError(t, err)
		r := httptest.NewRequest("GET", "/my-test", nil)
		handler := servedirectory.Handle("/my-test", dir)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		assert.Equal(t, "start here", w.Body.String())
	})

	t.Run("returns 404", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/static/does-not-exist.html", nil)
		handler := servedirectory.Handle("/static", dir)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
