package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registryCmd represents the registry command
var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Get files, objects, and work items from the APTrust Registry",
	Long: `Get files, objects, and work items from the APTrust Registry.
	Full online documentation:

      https://aptrust.github.io/userguide/partner_tools/
		
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Retrieve data from the APTrust registry. See subcommands for more info.")
	},
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
