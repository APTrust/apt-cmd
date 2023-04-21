package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Full online documentation:

https://aptrust.github.io/userguide/partner_tools/

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Retrieve data from the APTrust registry. See subcommands for more info.")
	},
}

func init() {
	registryCmd.AddCommand(getCmd)
}
