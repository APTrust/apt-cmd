package cmd

import (
	"fmt"
	"os"

	"github.com/APTrust/dart-runner/bagit"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a BagIt bag",
	Long:  `Create a BagIt bag`,
	Run: func(cmd *cobra.Command, args []string) {
		profileName := cmd.Flag("profile").Value.String()
		pathToDir := ""
		if len(args) > 0 {
			pathToDir = args[0]
		}
		if profileName == "" || pathToDir == "" {
			fmt.Println("Profile and path to directory are required.")
			os.Exit(EXIT_USER_ERR)
		}
		profile, err := LoadProfile(profileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}

		tags := GetTagValues(args)

		fmt.Println("Profile Name:", profileName)
		fmt.Println("Directory:", pathToDir)
		fmt.Println("Profile:", profile.Name)
		fmt.Println("Tag Values:")
		for _, t := range tags {
			fmt.Println("File:", t.TagFile, "Name:", t.TagName, "Value:", t.UserValue)
		}
	},
}

func init() {
	bagCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("profile", "p", "", "BagIt profile: 'aptrust', 'btr' or 'empty'")
}

func validateTags(profile *bagit.Profile, tags []*bagit.TagDefinition) {

	// TODO: Make sure that we got all required tags, and that
	//       tag values are valid. For bagit.txt, default to
	//       v1.0 and UTF-8 if no values are supplied.

}
