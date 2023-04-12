package cmd

import (
	"fmt"
	"os"
	"strings"

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

	// TODO: Add flag to specify manifest algorithms?
}

func EnsureDefaultTags(tags []*bagit.TagDefinition) {
	bagitVersion := FindTag(tags, "bagit.txt", "BagIt-Version")
	if bagitVersion == nil {
		versionTag := &bagit.TagDefinition{
			TagFile:   "bagit.txt",
			TagName:   "BagIt-Version",
			UserValue: "1.0",
		}
		tags = append(tags, versionTag)
	}
	encoding := FindTag(tags, "bagit.txt", "Tag-File-Character-Encoding")
	if encoding == nil {
		encodingTag := &bagit.TagDefinition{
			TagFile:   "bagit.txt",
			TagName:   "Tag-File-Character-Encoding",
			UserValue: "UTF-8",
		}
		tags = append(tags, encodingTag)
	}
}

func ValidateTags(profile *bagit.Profile, tags []*bagit.TagDefinition) []string {
	errors := make([]string, 0)
	for _, tagDef := range profile.Tags {
		hasValue := false
		userTag := FindTag(tags, tagDef.TagFile, tagDef.TagName)
		if tagDef.Required && userTag == nil {
			errors = append(errors, fmt.Sprintf("Required tag %s/%s is missing.", tagDef.TagFile, tagDef.TagName))
			continue
		}
		if userTag != nil && userTag.UserValue != "" {
			hasValue = true
		}
		if !tagDef.IsLegalValue(userTag.UserValue) {
			errors = append(errors, fmt.Sprintf("Tag %s/%s assigned illegal value '%s'. Valid values are: %s.", tagDef.TagFile, tagDef.TagName, userTag.UserValue, strings.Join(tagDef.Values, ",")))
			continue
		}
		if tagDef.Required && !tagDef.EmptyOK && !hasValue {
			errors = append(errors, fmt.Sprintf("Tag %s/%s is present but value cannot be empty. Please assign a value.", tagDef.TagFile, tagDef.TagName))
		}
	}
	return errors
}

// TODO: Change this to find tags? Tags can repeat.
func FindTag(tags []*bagit.TagDefinition, tagFile, tagName string) *bagit.TagDefinition {
	for _, tag := range tags {
		if tag.TagFile == tagFile && tag.TagName == tagName {
			return tag
		}
	}
	return nil
}
