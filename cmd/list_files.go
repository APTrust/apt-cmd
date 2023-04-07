package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "List files from the APTrust Registry",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		EnsureDefaultListParams(urlValues)
		resp := client.GenericFileList(urlValues)
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
	},
}

func init() {
	listCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
