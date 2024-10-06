package config

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/http-everything/httpe/pkg/share/timeunit"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DefaultServerAddress  = "0.0.0.0:3000"
	DefaultDataDir        = "./"
	DefaultDataRetention  = "1d"
	DefaultLogLevel       = "info"
	DefaultConfigFilename = "httpe.conf"
	EnvPrefix             = "httpe"
)

var (
	ErrAddressIncludesScheme  = errors.New("server address must not include scheme")
	ErrUnableToAccessCertFile = errors.New("failed to open/access cert_file")
	ErrUnableToAccessKeyFile  = errors.New("failed to open/access key_file")
	ErrCertOrKeyMissing       = errors.New("to activate TLS you must provide cert AND key")
	ErrNoRulesFile            = errors.New("no rules file specified")
	ErrRulesFileNotReadable   = errors.New("rules file not found or not readable")
	ErrBadSMTPServer          = errors.New("SMTP server is not a valid hostname or IP address")
)

// SvrConfig represents the config settings for the server
type SvrConfig struct {
	Address       string `mapstructure:"address"`
	DataDir       string `mapstructure:"data_dir"`
	DataRetention string `mapstructure:"data_retention"`
	CertFile      string `mapstructure:"cert_file"`
	KeyFile       string `mapstructure:"key_file"`
	AccessLogFile string `mapstructure:"access_log_file"`
	LogFile       string `mapstructure:"log_file"`
	LogLevel      string `mapstructure:"log_level"`
	RulesFile     string `mapstructure:"rules_file"`
	ValidateOnly  bool   `mapstructure:"validate"`
	DumpRules     bool   `mapstructure:"dump_rules"`
}

type SMTPConfig struct {
	Server   string `mapstructure:"server"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// Config is used for managing the license server config values
type Config struct {
	S    *SvrConfig  `mapstructure:"server"`
	SMTP *SMTPConfig `mapstructure:"smtp"`

	pFlags *pflag.FlagSet
	v      *viper.Viper
}

// New returns a new Config with the specified pFlags
func New(pFlags *pflag.FlagSet) (cfg *Config) {
	cfg = &Config{
		pFlags: pFlags,
	}
	return cfg
}

// LoadAndValidate loads and validates configuration values from either the path or reader specified
func (c *Config) LoadAndValidate(cfgPath *string, cfgReader io.Reader) (err error) {
	c.Setup()

	_, err = c.Load(cfgPath, cfgReader)
	if err != nil {
		return err
	}

	err = c.Validate()
	if err != nil {
		return err
	}

	return nil
}

// Setup configures viper for reading the config values
func (c *Config) Setup() {
	viperCfg := viper.New()
	viperCfg.SetConfigType("toml")

	viperCfg.SetDefault("server.address", DefaultServerAddress)
	viperCfg.SetDefault("server.data_dir", DefaultDataDir)
	viperCfg.SetDefault("server.data_retention", DefaultDataRetention)
	viperCfg.SetDefault("server.log_level", DefaultLogLevel)

	if c.pFlags != nil {
		_ = viperCfg.BindPFlag("server.address", c.pFlags.Lookup("address"))
		_ = viperCfg.BindPFlag("server.data_dir", c.pFlags.Lookup("data-dir"))
		_ = viperCfg.BindPFlag("server.data_retention", c.pFlags.Lookup("data-retention"))
		_ = viperCfg.BindPFlag("server.cert_file", c.pFlags.Lookup("cert-file"))
		_ = viperCfg.BindPFlag("server.key_file", c.pFlags.Lookup("key-file"))
		_ = viperCfg.BindPFlag("server.access_log_file", c.pFlags.Lookup("access-log-file"))
		_ = viperCfg.BindPFlag("server.log_file", c.pFlags.Lookup("log-file"))
		_ = viperCfg.BindPFlag("server.log_level", c.pFlags.Lookup("log-level"))
		_ = viperCfg.BindPFlag("server.rules_file", c.pFlags.Lookup("rules-file"))
		_ = viperCfg.BindPFlag("server.dump_rules", c.pFlags.Lookup("dump-rules"))
		_ = viperCfg.BindPFlag("server.validate", c.pFlags.Lookup("validate"))
	}

	viperCfg.SetEnvPrefix(EnvPrefix)
	replacer := strings.NewReplacer(".", "_")
	viperCfg.SetEnvKeyReplacer(replacer)
	viperCfg.AutomaticEnv()
	c.v = viperCfg
}

// Load reads the config from the reader (if non nil) or a file (if the reader
// isn't specified). If not using the reader, then the file will either be
// the one specified via the config flag or the default config filename.
func (c *Config) Load(cfgPath *string, cfgReader io.Reader) (cfg *Config, err error) {
	if cfgReader == nil {
		if *cfgPath != "" {
			c.v.SetConfigFile(*cfgPath)
		} else {
			c.v.AddConfigPath(".")
			c.v.SetConfigName(DefaultConfigFilename)
		}

		if err := c.v.ReadInConfig(); err != nil {
			//nolint:errorlint // viper.ConfigFileNotFoundError is a type and not a value
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("unable to read config file: %w", err)
			}
		}
	} else {
		if err := c.v.ReadConfig(cfgReader); err != nil {
			return nil, fmt.Errorf("unable to read config contents: %w", err)
		}
	}

	if err := c.v.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("failed parsing config file: %w", err)
	}

	return cfg, nil
}

// Validate validates the loaded config
func (c *Config) Validate() (err error) {
	address := c.S.Address
	if strings.Contains(address, "http") {
		return ErrAddressIncludesScheme
	}

	address = "https://" + address

	_, err = url.Parse(address)
	if err != nil {
		return err
	}

	if (c.S.KeyFile != "" && c.S.CertFile == "") || (c.S.KeyFile == "" && c.S.CertFile != "") {
		return ErrCertOrKeyMissing
	}

	if c.S.CertFile != "" {
		if !available(c.S.CertFile) {
			return ErrUnableToAccessCertFile
		}
	}
	if c.S.KeyFile != "" {
		if !available(c.S.KeyFile) {
			return ErrUnableToAccessKeyFile
		}
	}

	if c.S.RulesFile == "" {
		return ErrNoRulesFile
	}

	if !available(c.S.RulesFile) {
		return ErrRulesFileNotReadable
	}

	if c.SMTP != nil {
		if !isValidHostOrIPAddress(c.SMTP.Server) {
			return ErrBadSMTPServer
		}
	}

	_, err = timeunit.ParseDuration(c.S.DataRetention)
	if err != nil {
		return err
	}

	return nil
}

// available returns whether the given file is available
func available(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isValidHostOrIPAddress(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		return true
	}
	if govalidator.IsDNSName(ip) {
		return true
	}
	return false
}
