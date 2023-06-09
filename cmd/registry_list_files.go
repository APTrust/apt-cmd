package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "List files from the APTrust Registry",
	Long: `List files from the APTrust Registry, with filters.

--------------
Basic Examples
--------------
	
List files belonging to object test.edu/my_bag, ordered by identifer:
	
  apt-cmd registry list files intellectual_object_identifier='test.edu/my_bag' sort='identifier'

List only the first 10 files from that same bag:
	
  apt-cmd registry list files intellectual_object_identifier='test.edu/my_bag' sort='identifier' per_page=10

List files created after April 6, 2023

  apt-cmd registry list files created_at__gteq='2023-04-06'

Full online documentation:

  https://aptrust.github.io/userguide/partner_tools/
  
	`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(config, args)
		EnsureDefaultListParams(urlValues)
		resp := client.GenericFileList(urlValues)
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
	},
}

func init() {
	listCmd.AddCommand(filesCmd)
}
