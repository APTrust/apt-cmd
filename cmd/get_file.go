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
		urlValues := GetUrlValues(args)
		client, err := NewRegistryClient(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting Registry client:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		var resp *network.RegistryResponse
		fileID, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		if fileID > 0 {
			resp = client.GenericFileByID(fileID)
		} else {
			resp = client.GenericFileByIdentifier(urlValues.Get("identifier"))
		}
		//if resp.Error != nil {
		//	fmt.Fprintln(os.Stderr, "Error getting file from Registry:", resp.Error)
		//	os.Exit(EXIT_RUNTIME_ERR)
		//}
		data, _ := resp.RawResponseData()
		//if err != nil {
		//	fmt.Fprintln(os.Stderr, "Error reading response data:", err)
		//	os.Exit(EXIT_RUNTIME_ERR)
		//}
		fmt.Println(string(data))
		os.Exit(EXIT_OK)
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
