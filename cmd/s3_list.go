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

// s3ListCmd represents the list bucket command
var s3ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List items in an S3 bucket",
	Long: `Download a file from any S3 storage. For this to work,
you will need to have APTRUST_AWS_KEY and APTRUST_AWS_SECRET set in your 
environment, or in a config file specified with the --config flag.
	
Examples:

List items in my_bucket with prefix "photo":

    aptrust s3 list --host=s3.amazonaws.com --bucket=my_bucket --prefix=photo

List 10 items in my_bucket with prefix "photo", using plain text output:

    aptrust s3 list --host=s3.amazonaws.com \
                    --bucket=my_bucket \
                    --prefix=photo \
                    --maxitems=10 \
                    --format=text

	`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket := cmd.Flags().Lookup("bucket").Value.String()
		if bucket == "" {
			fmt.Fprintln(os.Stderr, "Missing required param --bucket")
			os.Exit(EXIT_USER_ERR)
		}
		s3Host := cmd.Flags().Lookup("host").Value.String()
		if s3Host == "" {
			fmt.Fprintln(os.Stderr, "Missing required param --host")
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
		logger.Infof("Listing up to %d items from %s/%s with prefix '%s'", maxKeys, s3Host, bucket, prefix)
		client := NewS3Client(config, s3Host)

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
	s3Cmd.AddCommand(s3ListCmd)
	s3ListCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3ListCmd.Flags().StringP("bucket", "b", "", "Bucket to list")
	s3ListCmd.Flags().StringP("prefix", "p", "", "List objects with this prefix")
	s3ListCmd.Flags().IntP("maxitems", "m", 50, "Maximum number of items to list (default = 50)")
	s3ListCmd.Flags().StringP("format", "f", "", "Output format: 'text' or 'json' (default = 'json')")
}
