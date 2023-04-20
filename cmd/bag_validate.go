package cmd

import (
	"embed"
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

  aptrust bag validate my_bag.tar
  aptrust bag validate -p aptrust my_bag.tar

To validate a bag using the Beyond the Repository (BTR) profile:

  aptrust bag validate -p btr my_bag.tar

To validate a bag using the empty profile:

  aptrust bag validate -p empty my_bag.tar

The empty profile simply ensures the bag is valid according to the general
BagIt specification. 

Limitations:

The validator only works with tarred bags and will not validate fetch.txt files.

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
		profile, err := LoadProfile(profileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
		logger.Debugf("Validating bag %s using profile %s", pathToBag, profile.Name)
		validator, err := bagit.NewValidator(pathToBag, profile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Can't create validator.", err.Error())
			os.Exit(EXIT_RUNTIME_ERR)
		}
		err = validator.ScanBag()
		if err != nil {
			fmt.Println("Bag is invalid due to the following errors:")
			fmt.Println(err.Error())
			os.Exit(EXIT_BAG_INVALID)
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
	bagCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("profile", "p", "", "BagIt profile: 'aptrust', 'btr' or 'empty'")
}
