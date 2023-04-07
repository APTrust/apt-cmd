package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// objectsCmd represents the objects command
var objectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "List object records from the APTrust Registry.",
	Long: `List objects from the APTrust Registry, with filters.

	--------------
	Basic Examples
	--------------
		
	List 20 objects ordered by identifer:
		
		aptrust list objects sort='identifier' per_page='20'
	
	List 20 objects reverse ordered by identifer:
		
		aptrust list objects sort='identifier__desc' per_page='20'
	
	List objects created after April 6, 2023
	
		aptrust list files created_at__gteq='2023-04-06'
	
		`,
	Run: func(cmd *cobra.Command, args []string) {
		client, urlValues := InitRegistryRequest(args)
		EnsureDefaultListParams(urlValues)
		resp := client.IntellectualObjectList(urlValues)
		data, _ := resp.RawResponseData()
		PrettyPrintJSON(data)
		os.Exit(EXIT_OK)
	},
}

func init() {
	listCmd.AddCommand(objectsCmd)
}
