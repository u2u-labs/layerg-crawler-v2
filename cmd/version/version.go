package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	BuildTime = ""
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", Version)
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Build Time:", BuildTime)
	},
}
