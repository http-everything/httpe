package cmd

import (
	"fmt"
	"os"

	"http-everything/httpe/pkg/config"

	"github.com/spf13/cobra"
)

var (
	RootCmd *cobra.Command
	CfgPath *string
)

func Execute() {
	RootCmd = &cobra.Command{
		Use:   "httpe",
		Short: "Start the httpe server",
		Long:  "The HTTPE server allows you to trigger a variety of actions via HTTP requests.",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	pFlags := RootCmd.PersistentFlags()

	// note this isn't currently working due to bug with cobra but leave in case bug fixed
	// https://github.com/spf13/cobra/issues/1033
	RootCmd.Flags().SortFlags = false

	CfgPath = pFlags.StringP("config", "c", "", "specify the config file to use")

	pFlags.StringP("address", "a", config.DefaultServerAddress, "set the listen address for the server")
	pFlags.StringP("data-dir", "d", config.DefaultDataDir, "set the data directory")
	pFlags.String("log-level", config.DefaultLogLevel, "specify server log level. either error, info, or debug.")
	pFlags.StringP("log-file", "l", "", "specify server log file")
	pFlags.StringP("rules-file", "r", "", "specify rules to map route to actions")
	pFlags.String("access-log-file", "", "set the access log file in the apache common log format. use '-' for writing to stdout.")
	pFlags.String("cert-file", "", "specify the TLS certificate file")
	pFlags.String("key-file", "", "specify the TLS key file")
	pFlags.Bool("validate", false, "validate configuration and rules, then exit")
	pFlags.BoolP("version", "v", false, "print version information")
	pFlags.Bool("dump-rules", false, "dump a json representation of the rules yaml, skips validation")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
