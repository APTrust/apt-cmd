package cmd

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"path"

	"github.com/APTrust/dart-runner/constants"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

func init() {
	s3Cmd.AddCommand(s3downloadCmd)
	s3downloadCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3downloadCmd.Flags().StringP("bucket", "b", "", "Bucket to download from")
	s3downloadCmd.Flags().StringP("key", "k", "", "Key (name of object) to download")
	s3downloadCmd.Flags().StringP("save-as", "s", "", "Name the file in which to save the download")
}

const (
	minChunkSize    = 64 * 1024 * 1024       // 64 MiB
	maxChunkSize    = 5 * 1024 * 1024 * 1024 // 5 GiB (S3 max part size)
	targetPartCount = 10000
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

    apt-cmd s3 download --host=s3.amazonaws.com --bucket="my-bucket" --key='photo_001.jpg'

Download the same file and save it with a custom name on your desktop:

    apt-cmd s3 download --host=s3.amazonaws.com  \
               --bucket="my-bucket" \
               --key='photo_001.jpg' \
               --save-as="$HOME/Desktop/vacation.jpg"

Full online documentation:

  https://aptrust.github.io/userguide/partner_tools/

`,
	Run: func(cmd *cobra.Command, args []string) {
		config.ValidateAWSCredentials()

		s3Host := GetFlagValue(cmd.Flags(), "host", "Missing required param --host")
		bucket := GetFlagValue(cmd.Flags(), "bucket", "Missing required param --bucket")
		key := GetFlagValue(cmd.Flags(), "key", "Missing required param --key")

		saveas := cmd.Flags().Lookup("save-as").Value.String()
		if saveas == "" {
			saveas = key
		}
		_stat, _ := os.Stat(saveas)
		if _stat != nil && _stat.IsDir() {
			saveas = path.Join(saveas, key)
		}
		logger.Debugf("Downloading object %s from %s/%s", key, s3Host, bucket)
		client := NewS3Client(config, s3Host)

		info, err := client.StatObject(context.Background(), bucket, key, minio.StatObjectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting size of S3 object:", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		var readCloser io.ReadCloser
		if info.Size > constants.MaxS3RequestSize {
			logger.Debugf("Using GetLargeObject to get %s from %s/%s", key, s3Host, bucket)
			readCloser, err = GetLargeObject(client, bucket, key, info.Size)
		} else {
			logger.Debugf("Using GetObject to get %s from %s/%s", key, s3Host, bucket)
			readCloser, err = client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error retrieving S3 object:", err)
			os.Exit(EXIT_REQUEST_ERROR)
		}
		defer readCloser.Close()
		outfile, err := os.Create(saveas)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening output file:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		_, err = io.Copy(outfile, readCloser)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing output file:", err)
			os.Exit(EXIT_RUNTIME_ERR)
		}
		fmt.Printf(`{ "result": "OK", "message": "S3 object %s saved to file %s" }`, key, saveas)
		fmt.Println("")
		os.Exit(EXIT_OK)
	},
}

// ComputeChunkSize computes the chunk size for retrieving a large
// object (> 5TB) from S3.
func ComputeChunkSize(objectSize int64) int64 {
	chunkSize := int64(math.Ceil(float64(objectSize / targetPartCount)))
	if chunkSize < minChunkSize {
		chunkSize = minChunkSize
	}
	if chunkSize > maxChunkSize {
		chunkSize = maxChunkSize
	}
	return chunkSize
}

// GetLargeObject returns a ReadCloser to download large objects
// (> 5TB) from S3.
func GetLargeObject(client *minio.Client, bucket, key string, size int64) (io.ReadCloser, error) {
	chunkSize := ComputeChunkSize(size)
	pr, pw := io.Pipe()

	go func() {
		var writeErr error
		var offset int64

		for offset < size {
			end := offset + chunkSize - 1
			if end >= size {
				end = size - 1
			}

			opts := minio.GetObjectOptions{}
			if err := opts.SetRange(offset, end); err != nil {
				writeErr = err
				break
			}

			obj, err := client.GetObject(context.Background(), bucket, key, opts)
			if err != nil {
				writeErr = err
				break
			}

			if _, err := io.Copy(pw, obj); err != nil {
				obj.Close()
				writeErr = err
				break
			}
			obj.Close()

			offset = end + 1
		}

		pw.CloseWithError(writeErr)
	}()

	return pr, nil
}
