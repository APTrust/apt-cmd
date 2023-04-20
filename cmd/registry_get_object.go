package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/APTrust/preservation-services/network"
	"github.com/spf13/cobra"
)

// objectCmd represents the object command
var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "Retrieve object metadata from the APTrust Registry",
	Long: `Retrieve a JSON record from the APTrust registry describing an
intellectual object. Object identifiers are strings,
such as 'example.edu/photos'. Ids are numeric.

Examples:

apt-cmd registry get object <object_identifier>
apt-cmd registry get object <object_id>

Full online documentation:

https://aptrust.github.io/userguide/partner_tools/

`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(config, args)
		var resp *network.RegistryResponse
		id, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		identifier := urlValues.Get("identifier")
		if id > 0 {
			resp = client.IntellectualObjectByID(id)
		} else if identifier != "" {
			resp = client.IntellectualObjectByIdentifier(identifier)
		} else {
			fmt.Fprintln(os.Stderr, "This call requires either an id or an identifier")
			os.Exit(EXIT_USER_ERR)
		}
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
	},
}

func init() {
	getCmd.AddCommand(objectCmd)
}
