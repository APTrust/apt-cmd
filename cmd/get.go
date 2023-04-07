package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a record from the APTrust Registry",
	Long: `Retrieve a single record from the APTrust Registry.
Records are in json format. Object and file identifiers are strings,
such as 'example.edu/photos' (object identifier) or
'example.edu/photos/data/image1.jpg' (file identifier). Ids are numeric.

Examples to retrieve individual records:

aptrust get object <object_identifier>
aptrust get object <object_id>

aptrust get file <file_identifier>
aptrust get file <file_id>

aptrust get workitem <item_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
