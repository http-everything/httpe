package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	HTTPEServerVersion    = "0.0-src"
	HTTPEServerBuildTime  = ""
	HTTPEServerLocalBuild = ""
	HTTPEServerGitRef     = ""
	HTTPEServerCommitID   = ""
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of HTTPE",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("HTTPE version %s\n", HTTPEServerVersion)
	},
}
