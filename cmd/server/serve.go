package cmd

import (
	"context"
	"fmt"
	"http-everything/httpe/pkg/rules"
	"io"
	"os"

	"http-everything/httpe/pkg/config"
	"http-everything/httpe/pkg/server"
	"http-everything/httpe/pkg/share/logger"
	"http-everything/httpe/pkg/share/version"
)

const (
	accessLogPermissions = 0644
)

// serve starts and runs the HTTPE server
func serve() {
	pFlags := RootCmd.PersistentFlags()

	pv, _ := pFlags.GetBool("version")
	if pv {
		fmt.Printf("httpe %s\n", version.HTTPEServerVersion)
		return
	}

	cfg := config.New(pFlags)

	err := cfg.LoadAndValidate(CfgPath, nil)
	if err != nil {
		reportErrorAndExit(nil, fmt.Errorf("unable to load and validate config: %w", err))
		return
	}

	baseLogger, accessLogWriter, err := setupLogs(cfg)
	if err != nil {
		reportErrorAndExit(nil, fmt.Errorf("unable to open log files: %w", err))
		reportErrorAndExit(nil, err)
		return
	}

	baseLogger.Infof("running version %s", version.HTTPEServerVersion)
	if *CfgPath != "" {
		baseLogger.Infof("config loaded from %s", *CfgPath)
	} else {
		baseLogger.Infof("config loaded from %s", config.DefaultConfigFilename)
	}

	rulesCfg, err := rules.Read(cfg.S.RulesFile, baseLogger.Fork("rules"))
	if err != nil {
		reportErrorAndExit(baseLogger, err)
		return
	}
	if cfg.S.DumpRules {
		fmt.Println(rules.YamlToJSON(cfg.S.RulesFile))
		return
	}
	err = rulesCfg.Validate()
	if err != nil {
		reportErrorAndExit(baseLogger, err)
	}
	if cfg.S.ValidateOnly {
		// End here if only the validation has been requested
		return
	}

	svr, err := server.New(cfg, rulesCfg.Rules, baseLogger, accessLogWriter)
	if err != nil {
		reportErrorAndExit(baseLogger, fmt.Errorf("unable to setup HTTPE server: %w", err))
		return
	}

	svr.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = svr.Serve(ctx, true)
	if err != nil {
		reportErrorAndExit(baseLogger, fmt.Errorf("unable to start HTTPE server: %w", err))
		return
	}

	svr.Shutdown()

	baseLogger.Shutdown()
}

// setupLogs sets up the regular and access logs using the config specified
func setupLogs(cfg *config.Config) (baseLogger *logger.Logger, accessLogWriter io.Writer, err error) {
	baseLogger, err = logger.New("serve", cfg.S.LogFile, cfg.S.LogLevel)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open server log: %w", err)
	}
	var AccessLogWriter = os.Stdout
	if cfg.S.AccessLogFile != "" {
		AccessLogWriter, err = os.OpenFile(cfg.S.AccessLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, accessLogPermissions)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open access log: %w", err)
		}
		baseLogger.Infof("Using access log file %s", cfg.S.AccessLogFile)
	}

	return baseLogger, AccessLogWriter, nil
}

// reportErrorAndExit is a simple helper fn to log to screen and the regular log
func reportErrorAndExit(l *logger.Logger, err error) {
	if l != nil {
		l.Errorf(err.Error())
	}
	fmt.Println("error: " + err.Error())
	os.Exit(1)
}
