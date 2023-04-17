/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

// s3deleteCmd represents the s3delete command
var s3deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an object from S3 storage",
	Long: `Delete an object from any S3-compatible service. For this to work,
you will need to have APTRUST_AWS_KEY and APTRUST_AWS_SECRET set in your 
environment, or in a config file specified with the --config flag.
	
Example:
	
Delete object photo.jpg from my-bucket on AWS S3:
	
    s3delete --host=s3.amazonaws.com --bucket="my-bucket" --key='photo.jpg' 

Note: This returns exit status zero and '{ "result": "OK" }' if the key is 
successfully deleted or if the key wasn't in the bucket to begin with.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ValidateAWSCredentials()

		s3Host := GetFlagValue(cmd.Flags(), "host", "Missing required param --host")
		bucket := GetFlagValue(cmd.Flags(), "bucket", "Missing required param --bucket")
		key := GetFlagValue(cmd.Flags(), "key", "Missing required param --key")
		if LooksLikePreservationBucket(bucket) {
			fmt.Fprintln(os.Stderr, "Deletion from preservation bucket not allowed")
			os.Exit(EXIT_USER_ERR)
		}

		logger.Debugf("Deleting object %s from %s/%s", key, s3Host, bucket)
		client := NewS3Client(config, s3Host)
		err := client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting object: ", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		fmt.Printf(`{ "result": "OK", "message": "Deleted %s/%s" }`, bucket, key)
		fmt.Println("")
		os.Exit(EXIT_OK)
	},
}

func init() {
	s3Cmd.AddCommand(s3deleteCmd)
	s3deleteCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3deleteCmd.Flags().StringP("bucket", "b", "", "Bucket containing object to delete")
	s3deleteCmd.Flags().StringP("key", "k", "", "Key (name of object) to delete")
}
