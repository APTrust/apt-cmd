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

apt-cmd registry get workitem <id>

Full online documentation:

https://aptrust.github.io/userguide/partner_tools/

`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(config, args)
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
}
