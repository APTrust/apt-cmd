package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files, objects, or work items from the APTrust Registry",
	Long: `List files, objects, or work items from the APTrust Registry.
	Full online documentation:

	  https://aptrust.github.io/userguide/partner_tools/		
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See subcomands.")
	},
}

func init() {
	registryCmd.AddCommand(listCmd)
}
