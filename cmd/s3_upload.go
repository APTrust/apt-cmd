package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

// s3uploadCmd represents the s3upload command
var s3uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to an S3-compatible service",
	Long: `Upload a file to any S3-compatible service. For this to work,
you will need to have APTRUST_AWS_KEY and APTRUST_AWS_SECRET set in your 
environment, or in a config file specified with the --config flag.
	
Examples:
	
Upload file photo.jpg to Amazon's S3 service:
	
    s3upload --host=s3.amazonaws.com --bucket="my-bucket" photo.jpg 
	
Upload the same file, but call it renamed.jpg in S3:
	
    s3upload --host=s3.amazonaws.com  \
             --bucket="my-bucket" \
             --key='renamed.jpg' \
             photo.jpg
			 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ValidateAWSCredentials()
		file := ""
		if len(args) > 0 {
			file = args[0]
		}
		if file == "" {
			fmt.Fprintln(os.Stderr, "Missing required arg file")
			os.Exit(EXIT_USER_ERR)
		}
		fstat, err := os.Stat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "File", file, "is missing or unreadable")
			os.Exit(EXIT_USER_ERR)
		}
		if fstat.IsDir() {
			fmt.Fprintln(os.Stderr, "File", file, "is a directory")
			os.Exit(EXIT_USER_ERR)
		}
		s3Host := GetFlagValue(cmd.Flags(), "host", "Missing required param --host")
		bucket := GetFlagValue(cmd.Flags(), "bucket", "Missing required param --bucket")

		if LooksLikePreservationBucket(bucket) {
			fmt.Fprintln(os.Stderr, "Upload to preservation bucket not allowed")
			os.Exit(EXIT_USER_ERR)
		}

		key := cmd.Flags().Lookup("key").Value.String()
		if key == "" {
			key = path.Base(file)
		}

		logger.Infof("Uploading file %s to %s/%s/%s", file, s3Host, bucket, key)
		client := NewS3Client(config, s3Host)
		uploadInfo, err := client.FPutObject(context.Background(), bucket, key, file, minio.PutObjectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error uploading file:", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		data, err := json.MarshalIndent(uploadInfo, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error serializing JSON response from S3 server:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		fmt.Println(string(data))
		os.Exit(EXIT_OK)

	},
}

func init() {
	s3Cmd.AddCommand(s3uploadCmd)
	s3uploadCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3uploadCmd.Flags().StringP("bucket", "b", "", "Bucket to upload from")
	s3uploadCmd.Flags().StringP("key", "k", "", "Key (name of object) to download")
}
