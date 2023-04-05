package cmd

import (
	"fmt"

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
		fmt.Println("get file called")
		fmt.Println(args)
		fmt.Println(ParseArgPairs(args))
		fmt.Println(GetUrlValues(args))
	},
}

func init() {
	getCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
