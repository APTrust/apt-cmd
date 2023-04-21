package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var versionNumber = "v3.0.0-beta"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info and exit",
	Long:  `Print version info and exit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apt-cmd (APTrust partner tools)")
		fmt.Println(versionNumber, runtime.GOOS, runtime.GOARCH)
		os.Exit(EXIT_OK)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
