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

aptrust get file <file_identifier>
aptrust get file <file_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		var resp *network.RegistryResponse
		fileID, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		if fileID > 0 {
			resp = client.GenericFileByID(fileID)
		} else {
			resp = client.GenericFileByIdentifier(urlValues.Get("identifier"))
		}
		data, _ := resp.RawResponseData()
		fmt.Println(string(data))
		os.Exit(EXIT_OK)
	},
}

func init() {
	getCmd.AddCommand(fileCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
