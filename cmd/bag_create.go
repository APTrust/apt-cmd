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
	createCmd.Flags().StringSliceP("manifest-algs", "m", []string{"sha256"}, "Manifest algorithms. Use comma-separated list for multiple. Supported algorithms: md5, sha1, sha256, sha512. Default is sha256.")
}

func EnsureDefaultTags(tags []*bagit.TagDefinition) []*bagit.TagDefinition {
	bagitVersion := FindTag(tags, "bagit.txt", "BagIt-Version")
	if bagitVersion == nil {
		versionTag := &bagit.TagDefinition{
			TagFile:   "bagit.txt",
			TagName:   "BagIt-Version",
			UserValue: "1.0",
		}
		tags = append(tags, versionTag)
	} else if bagitVersion.GetValue() == "" {
		bagitVersion.UserValue = "1.0"
	}
	encoding := FindTag(tags, "bagit.txt", "Tag-File-Character-Encoding")
	if encoding == nil {
		encodingTag := &bagit.TagDefinition{
			TagFile:   "bagit.txt",
			TagName:   "Tag-File-Character-Encoding",
			UserValue: "UTF-8",
		}
		tags = append(tags, encodingTag)
	} else if encoding.GetValue() == "" {
		encoding.UserValue = "UTF-8"
	}
	return tags
}

// ValidateTags verifies that tags required by the BagIt profile are
// present and contain valid values. We check this BEFORE bagging because
// in case where the user is packaging 500+ GB, they don't want to wait
// two hours to find out their bag is invalid.
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
		if userTag != nil && !tagDef.IsLegalValue(userTag.UserValue) {
			errors = append(errors, fmt.Sprintf("Tag %s/%s assigned illegal value '%s'. Valid values are: %s.", tagDef.TagFile, tagDef.TagName, userTag.UserValue, strings.Join(tagDef.Values, ",")))
			continue
		}
		if tagDef.Required && !tagDef.EmptyOK && !hasValue {
			errors = append(errors, fmt.Sprintf("Tag %s/%s is present but value cannot be empty. Please assign a value.", tagDef.TagFile, tagDef.TagName))
		}
	}
	return errors
}

// ValidateManifestAlgorithms checks to see whether the user-specified manifest
// algorithms are allowed by the profile, and whether the user specified all
// of the profile's required algorithms. We do this work up front, before creating
// the bag, to avoid creating an invalid bag.
func ValidateManifestAlgorithms(profile *bagit.Profile, algs []string) []string {
	errors := make([]string, 0)
	for _, alg := range algs {
		isAllowed := false
		for _, allowedAlg := range profile.ManifestsAllowed {
			if allowedAlg == alg {
				isAllowed = true
			}
		}
		if !isAllowed {
			errors = append(errors, fmt.Sprintf("Manifest algorithm '%s' is not allowed in profile %s.", alg, profile.Name))
		}
	}
	for _, requiredAlg := range profile.ManifestsRequired {
		foundRequiredAlg := false
		for _, alg := range algs {
			if alg == requiredAlg {
				foundRequiredAlg = true
			}
		}
		if !foundRequiredAlg {
			errors = append(errors, fmt.Sprintf("Profile %s requires manifest algorithm %s", profile.Name, requiredAlg))
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
