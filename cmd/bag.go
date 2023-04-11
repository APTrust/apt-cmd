package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bagCmd represents the bag command
var bagCmd = &cobra.Command{
	Use:   "bag",
	Short: "Create and validate BagIt bags.",
	Long:  `Create and validate BagIt bags.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("See subcommands")
	},
}

func init() {
	rootCmd.AddCommand(bagCmd)
}
