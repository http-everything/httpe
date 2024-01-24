package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"http-everything/httpe/pkg/config"
)

const (
	Address         = "127.0.0.0:3000"
	CertFile        = "../../testdata/certs/testcert.pem"
	KeyFile         = "../../testdata/certs/testkey.pem"
	NonExistingFile = "/tmp/aiWa0weshie4Shahcoh4"
)

func TestShouldValidateConfig(t *testing.T) {
	validServerConfig := &config.SvrConfig{
		Address:   Address,
		CertFile:  CertFile,
		KeyFile:   KeyFile,
		RulesFile: "../../testdata/rules/good/all.yaml",
	}
	cases := []struct {
		name      string
		cfg       *config.Config
		wantError error
	}{
		{
			name: "valid config",
			cfg: &config.Config{
				S: validServerConfig,
			},
			wantError: nil,
		},
		{
			name: "bad address (includes scheme)",
			cfg: &config.Config{
				S: &config.SvrConfig{
					Address:  "http://" + Address,
					CertFile: CertFile,
					KeyFile:  KeyFile,
				},
			},
			wantError: config.ErrAddressIncludesScheme,
		},
		{
			name: "cert file inaccessible",
			cfg: &config.Config{
				S: &config.SvrConfig{
					Address:  Address,
					CertFile: NonExistingFile,
					KeyFile:  KeyFile,
				},
			},
			wantError: config.ErrUnableToAccessCertFile,
		},
		{
			name: "key file inaccessible",
			cfg: &config.Config{
				S: &config.SvrConfig{
					Address:  Address,
					CertFile: CertFile,
					KeyFile:  NonExistingFile,
				},
			},
			wantError: config.ErrUnableToAccessKeyFile,
		},
		{
			name: "missing key file",
			cfg: &config.Config{
				S: &config.SvrConfig{
					Address: Address,
					KeyFile: NonExistingFile,
				},
			},
			wantError: config.ErrCertOrKeyMissing,
		},
		{
			name: "missing cert file",
			cfg: &config.Config{
				S: &config.SvrConfig{
					Address:  Address,
					CertFile: NonExistingFile,
				},
			},
			wantError: config.ErrCertOrKeyMissing,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()

			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShouldVerifyDefaults(t *testing.T) {
	os.Setenv("HTTPE_SERVER_RULES_FILE", "../../testdata/rules/good/all.yaml")
	//@Todo: Set rules files via pflags
	cfg := config.New(nil)
	cfg.Setup()

	cfgPath := "../../example.httpe.conf"
	_, err := cfg.Load(&cfgPath, nil)
	require.NoError(t, err)

	err = cfg.Validate()
	require.NoError(t, err)

	assert.Equal(t, "0.0.0.0:3000", cfg.S.Address)
	assert.Equal(t, "/var/lib/httpe", cfg.S.DataDir)
	assert.Equal(t, "/var/log/httpe/access.log", cfg.S.AccessLogFile)
	assert.Equal(t, "/var/log/httpe/server.log", cfg.S.LogFile)
	assert.Equal(t, "info", cfg.S.LogLevel)
}

func TestShouldVerifyEnvVars(t *testing.T) {
	os.Setenv("HTTPE_SERVER_ADDRESS", "127.0.0.1:3001")
	os.Setenv("HTTPE_SERVER_RULES_FILE", "../../testdata/rules/good/all.yaml")
	cfg := config.New(nil)
	cfg.Setup()

	cfgPath := "../../example.httpe.conf"
	_, err := cfg.Load(&cfgPath, nil)
	require.NoError(t, err)

	err = cfg.Validate()
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1:3001", cfg.S.Address)
}
