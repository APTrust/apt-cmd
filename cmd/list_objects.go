/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// objectsCmd represents the objects command
var objectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "List object records from the APTrust Registry.",
	Long:  ``,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// objectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// objectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
