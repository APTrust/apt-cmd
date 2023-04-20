package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// s3Cmd is the top-level command for s3 operations
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Upload, download, list and delete S3 objects",
	Long: `Upload, download, list and delete S3 objects.
For more info, run:

    apt-cmd s3 upload --help
    apt-cmd s3 download --help
    apt-cmd s3 list --help
    apt-cmd s3 delete --help

Full online documentation:

  https://aptrust.github.io/userguide/partner_tools/
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run aptrust s3 --help for more info")
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)
}
