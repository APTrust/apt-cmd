package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/APTrust/preservation-services/network"
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
		client, urlValues := InitRegistryRequest(args)
		var resp *network.RegistryResponse
		id, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		if id > 0 {
			resp = client.WorkItemByID(id)
		} else {
			fmt.Fprintln(os.Stderr, "This call requires an id (e.g. id=1234)")
			os.Exit(EXIT_USER_ERR)
		}
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
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
