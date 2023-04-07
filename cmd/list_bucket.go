/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

// bucketCmd represents the bucket command
var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket := cmd.Flags().Lookup("bucket").Value.String()
		if bucket == "" {
			fmt.Fprintln(os.Stderr, "Missing required param --bucket")
			os.Exit(EXIT_USER_ERR)
		}
		format := cmd.Flags().Lookup("format").Value.String()
		if format != "" && format != "text" && format != "json" {
			fmt.Fprintln(os.Stderr, "Unknown format:", format, ". Defaulting to json.")
		}
		prefix := cmd.Flags().Lookup("prefix").Value.String()
		maxKeys, err := strconv.Atoi(cmd.Flags().Lookup("maxitems").Value.String())
		if err != nil {
			maxKeys = 50
			logger.Info("Could not parse maxitems. Defaulting to 50.")
		}
		logger.Infof("Listing up to %d items from %s/%s with prefix '%s'", maxKeys, config.S3Host, bucket, prefix)
		client := GetS3Client()

		doneCh := make(chan struct{})
		defer close(doneCh)
		objectCh := client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: false,
			MaxKeys:   maxKeys,
		})

		objCount := 0
		for obj := range objectCh {
			objCount += 1
			if obj.Err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v", bucket, obj.Err)
				continue
			}
			if format == "text" {
				fmt.Println("Key:     ", obj.Key)
				fmt.Println("Etag:    ", obj.ETag)
				fmt.Println("Size:    ", humanize.Comma(obj.Size), "(", humanize.Bytes(uint64(obj.Size)), ")")
				fmt.Println("Modified:", obj.LastModified.Format(time.RFC3339))
				fmt.Println("----------------------------------------------------")
			} else {
				data, err := json.MarshalIndent(obj, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error serialized entry '%s' to JSON: %v", obj.Key, err)
					continue
				}
				fmt.Println(string(data))
			}
			if objCount >= maxKeys {
				// Note that Minio seems to take maxKeys as a suggestion
				// and doesn't always follow it. So we break manually when
				// we've hit our limit.
				// See https://github.com/minio/minio-go/issues/1536
				break
			}
		}
	},
}

func init() {
	listCmd.AddCommand(bucketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bucketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bucketCmd.Flags().StringP("bucket", "b", "", "Bucket to list")
	bucketCmd.Flags().StringP("prefix", "p", "", "List objects with this prefix")
	bucketCmd.Flags().IntP("maxitems", "m", 50, "Maximum number of items to list (default = 50)")
	bucketCmd.Flags().StringP("format", "f", "", "Output format: 'text' or 'json' (default = 'json')")
}
