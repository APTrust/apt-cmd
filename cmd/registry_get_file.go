package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/APTrust/preservation-services/network"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Retrieve file metadata from the APTrust Registry",
	Long: `Retrieve a JSON record from the APTrust registry describing a
generic file. File identifiers are strings,
such as 'example.edu/photos/data/image1.jpg'. Ids are numeric.

Examples:

aptrust registry get file <file_identifier>
aptrust registry get file <file_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(config, args)
		var resp *network.RegistryResponse
		id, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		identifier := urlValues.Get("identifier")
		if id > 0 {
			resp = client.GenericFileByID(id)
		} else if identifier != "" {
			resp = client.GenericFileByIdentifier(identifier)
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
	getCmd.AddCommand(fileCmd)
}
