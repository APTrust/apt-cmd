package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/APTrust/dart-runner/bagit"
	"github.com/spf13/cobra"
)

//go:embed profiles
var profiles embed.FS

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a bag using the APTrust, BTR, or empty BagIt profile.",
	Long: `Validate a bag according to a specific BagIt profile.
Currently, this supports only tarred bags. The following commands
validate a bag according to the APTrust BagIt profile:

aptrust validate my_bag.tar
aptrust validate -p aptrust my_bag.tar

To validate a bag using the Beyond the Repository (BTR) profile:

aptrust validate -p btr my_bag.tar

To validate a bag using the empty profile:

aptrust validate -p empty my_bag.tar

The empty profile simply ensures the bag is valid according to the general
BagIt specification.
`,
	Run: func(cmd *cobra.Command, args []string) {
		profileName := cmd.Flag("profile").Value.String()
		pathToBag := ""
		if len(args) > 0 {
			pathToBag = args[0]
		}
		if profileName == "" || pathToBag == "" {
			fmt.Println("Profile and path to bag are required.")
			os.Exit(EXIT_USER_ERR)
		}
		profile, err := loadProfile(profileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
		validator, err := bagit.NewValidator(pathToBag, profile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
		err = validator.ScanBag()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
		if validator.Validate() {
			fmt.Println("Bag is valid according to", profileName, "profile.")
			os.Exit(EXIT_OK)
		}
		fmt.Println("Bag is invalid due to the following errors:")
		for key, value := range validator.Errors {
			fmt.Println(key, ": ", value)
		}
		os.Exit(EXIT_BAG_INVALID)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("profile", "p", "", "BagIt profile: aptrust or btr")
}

func loadProfile(name string) (*bagit.Profile, error) {
	profile := &bagit.Profile{}
	var data []byte
	var err error
	switch name {
	case "aptrust":
		data, err = profiles.ReadFile("profiles/aptrust-v2.2.json")
	case "btr":
		data, err = profiles.ReadFile("profiles/btr-v1.0.json")
	case "empty":
		data, err = profiles.ReadFile("profiles/empty_profile.json")
	default:
		err = fmt.Errorf("missing or invalid profile. Only 'aptrust', 'btr' and 'empty' are supported")
	}
	if err == nil && len(data) > 1 {
		err = json.Unmarshal(data, profile)
	}
	return profile, err
}
