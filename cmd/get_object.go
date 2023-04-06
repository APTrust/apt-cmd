/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
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

aptrust get object <object_identifier>
aptrust get object <object_id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		var resp *network.RegistryResponse
		id, _ := strconv.ParseInt(urlValues.Get("id"), 10, 64)
		identifier := urlValues.Get("identifier")
		if id > 0 {
			resp = client.IntellectualObjectByID(id)
		} else if identifier != "" {
			resp = client.IntellectualObjectByIdentifier(identifier)
		} else {
			fmt.Fprintln(os.Stderr, "This call requires either an id or an identifier")
		}
		data, _ := resp.RawResponseData()
		fmt.Println(string(data))
		os.Exit(EXIT_OK)
	},
}

func init() {
	getCmd.AddCommand(objectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// objectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// objectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
