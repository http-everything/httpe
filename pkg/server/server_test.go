package server_test

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/http-everything/httpe/pkg/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/server"
	"github.com/http-everything/httpe/pkg/share/logger"
)

func TestShouldReturnNotFoundForInvalidRoute(t *testing.T) {
	testPath := "/some-path"
	cfg, testLogger := makeTestConfig(t)
	accessLog := bytes.Buffer{}
	ru := &[]rules.Rule{}

	svr, err := server.New(cfg, ru, testLogger, &accessLog)
	require.NoError(t, err)

	svr.Setup()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", testPath, nil)

	svr.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, accessLog.String(), fmt.Sprintf("GET %s HTTP/1.1\" 404 %d", testPath, 10))
}

func TestShouldConnectToTLSServer(t *testing.T) {
	cfg, testLogger := makeTestConfig(t)
	cfg.S.CertFile = "../../testdata/certs/testcert.pem"
	cfg.S.KeyFile = "../../testdata/certs/testkey.pem"
	ru := &[]rules.Rule{}

	svr, err := server.New(cfg, ru, testLogger, nil)
	require.NoError(t, err)

	svr.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = svr.Serve(ctx, false)
	require.NoError(t, err)
	// allow the TLS server a little time to initialise
	time.Sleep(200 * time.Millisecond)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//nolint:gosec
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	res, err := client.Get("https://127.0.0.1:3000/")
	require.NoError(t, err)
	defer res.Body.Close()
	svr.Shutdown()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestShouldConnect(t *testing.T) {
	cfg, testLogger := makeTestConfig(t)
	ru := &[]rules.Rule{
		{
			On: &rules.On{
				Path:    "/",
				Methods: []string{"get"},
			},
			AnswerContent: "hello world",
		},
	}

	svr, err := server.New(cfg, ru, testLogger, nil)
	require.NoError(t, err)

	svr.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = svr.Serve(ctx, false)
	require.NoError(t, err)
	// allow the server a little time to initialise
	time.Sleep(200 * time.Millisecond)

	client := &http.Client{}
	res, err := client.Get("http://127.0.0.1:3000/")
	require.NoError(t, err)
	defer res.Body.Close()
	svr.Shutdown()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "hello world", string(body))
}

func makeTestConfig(t *testing.T) (cfg *config.Config, l *logger.Logger) {
	t.Helper()

	l, _ = logger.New("test", "", "debug")

	cfg = &config.Config{
		S: &config.SvrConfig{
			Address: "127.0.0.1:3000",
		},
	}
	return cfg, l
}
