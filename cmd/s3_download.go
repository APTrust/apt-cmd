package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

// s3downloadCmd represents the s3download command
var s3downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from S3 storage",
	Long: `Download a file from any S3 storage. For this to work,
you will need to have APTRUST_AWS_KEY and APTRUST_AWS_SECRET set in your 
environment, or in a config file specified with the --config flag.

Examples:

Download a file from Amazon's S3 service into the current directory:

    s3download --host=s3.amazonaws.com --bucket="my-bucket" --key='photo_001.jpg' 

Download the same file and save it with a custom name on your desktop:

    s3download --host=s3.amazonaws.com  \
               --bucket="my-bucket" \
               --key='photo_001.jpg' \
               --save-as="$HOME/Desktop/vacation.jpg"
		   
`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ValidateAWSCredentials()

		s3Host := GetParam(cmd.Flags(), "host", "Missing required param --host")
		bucket := GetParam(cmd.Flags(), "bucket", "Missing required param --bucket")
		key := GetParam(cmd.Flags(), "key", "Missing required param --key")

		saveas := cmd.Flags().Lookup("save-as").Value.String()
		if saveas == "" {
			saveas = key
		}
		_stat, _ := os.Stat(saveas)
		if _stat != nil && _stat.IsDir() {
			saveas = path.Join(saveas, key)
		}
		logger.Infof("Downloading object %s from %s/%s", key, s3Host, bucket)
		client := GetS3Client(s3Host)
		obj, err := client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error retrieving S3 object:", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		defer obj.Close()
		outfile, err := os.Create(saveas)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening output file:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		_, err = io.Copy(outfile, obj)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing output file:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		os.Exit(EXIT_OK)
	},
}

func init() {
	s3Cmd.AddCommand(s3downloadCmd)
	s3downloadCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3downloadCmd.Flags().StringP("bucket", "b", "", "Bucket to download from")
	s3downloadCmd.Flags().StringP("key", "k", "", "Key (name of object) to download")
	s3downloadCmd.Flags().StringP("save-as", "s", "", "Name the file in which to save the download")
}
