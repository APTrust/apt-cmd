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
	Use:   "s3delete",
	Short: "Delete an object from S3 storage",
	Long: `Delete an object from any S3-compatible service. For this to work,
you will need to have APTRUST_AWS_KEY and APTRUST_AWS_SECRET set in your 
environment, or in a config file specified with the --config flag.
	
Example:
	
Download object photo.jpg from my-bucket on AWS S3:
	
	s3delete --host=s3.amazonaws.com --bucket="my-bucket" --key='photo.jpg' 

Note: This returns exit status zero and '{ "result": "OK" }' if the key is 
successfully deleted or if the key wasn't in the bucket to begin with.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ValidateAWSCredentials()

		s3Host := GetParam(cmd.Flags(), "host", "Missing required param --host")
		bucket := GetParam(cmd.Flags(), "bucket", "Missing required param --bucket")
		key := GetParam(cmd.Flags(), "key", "Missing required param --key")
		DisallowPreservationBucket(bucket)

		logger.Infof("Deleting object %s from %s/%s", key, s3Host, bucket)
		client := GetS3Client(s3Host)
		err := client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error deleting object: ", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		fmt.Println(`{ "result": "OK" }`)
		os.Exit(EXIT_OK)
	},
}

func init() {
	rootCmd.AddCommand(s3deleteCmd)
	s3deleteCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3deleteCmd.Flags().StringP("bucket", "b", "", "Bucket containing object to delete")
	s3deleteCmd.Flags().StringP("key", "k", "", "Key (name of object) to delete")
}
