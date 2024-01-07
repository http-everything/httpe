package cmd

import (
	"fmt"
	"http-everything/httpe/pkg/share/version"
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
			serve(false)
		},
	}

	RootCmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate config and rules",
		Long: "Validate config and rules then exit. On errors errors, exit code>0." +
			" Use in combination with '--log-level debug' to get all details of validation errors.",
		Run: func(cmd *cobra.Command, args []string) {
			serve(true)
		},
	})

	RootCmd.AddCommand(version.Cmd)

	pFlags := RootCmd.PersistentFlags()

	// note this isn't currently working due to bug with cobra but leave in case bug fixed
	// https://github.com/spf13/cobra/issues/1033
	RootCmd.Flags().SortFlags = false

	CfgPath = pFlags.StringP("config", "c", "", "specify the config file to use")

	pFlags.StringP("address", "a", config.DefaultServerAddress, "set the listen address for the server")
	pFlags.StringP("data-dir", "d", config.DefaultDataDir, "set the data directory for license json files")
	pFlags.String("log-level", config.DefaultLogLevel, "specify server log level. either error, info, or debug.")
	pFlags.StringP("log-file", "p", "", "specify server log file")
	pFlags.StringP("rules-file", "r", "", "specify rules to map route to actions")
	pFlags.String("access-log-file", "", "set the access log file in the apache common log format")
	pFlags.String("cert-file", "", "specify the TLS certificate file")
	pFlags.String("key-file", "", "specify the TLS key file")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
