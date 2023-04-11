package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registryCmd represents the registry command
var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Get files, objects, and work items from the APTrust Registry",
	Long:  `Get files, objects, and work items from the APTrust Registry`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See subcommands")
	},
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
