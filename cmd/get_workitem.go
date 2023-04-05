package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// workitemCmd represents the workitem command
var workitemCmd = &cobra.Command{
	Use:   "workitem",
	Short: "Retrieves a WorkItem record from the APTrust Registry",
	Long: `Retrieve a WorkItem record from the APTrust Registry. Use this
to check on the status of ingests, restorations and deletions. Id is a
number.

aptrust get workitem <id>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get workitem called")
	},
}

func init() {
	getCmd.AddCommand(workitemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workitemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workitemCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
