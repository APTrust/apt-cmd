package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// workitemsCmd represents the workitems command
var workitemsCmd = &cobra.Command{
	Use:   "workitems",
	Short: "List work item records from the APTrust registry.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		EnsureDefaultListParams(urlValues)
		resp := client.WorkItemList(urlValues)
		data, _ := resp.RawResponseData()
		fmt.Println(string(data))
		os.Exit(EXIT_OK)
	},
}

func init() {
	listCmd.AddCommand(workitemsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workitemsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// workitemsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
