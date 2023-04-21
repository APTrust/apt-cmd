package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var Version string
var CommitId string
var BuildDate string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info and exit",
	Long:  `Print version info and exit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("  apt-cmd (APTrust partner tools)")
		fmt.Printf("  Version %s on %s %s\n", Version, runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Build %s on %s\n", CommitId, BuildDate)
		os.Exit(EXIT_OK)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
