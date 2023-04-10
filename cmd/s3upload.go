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
	Use:   "s3upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		bucket := cmd.Flags().Lookup("bucket").Value.String()
		s3Host := cmd.Flags().Lookup("host").Value.String()
		if s3Host == "" {
			fmt.Fprintln(os.Stderr, "Missing required param --host")
			os.Exit(EXIT_USER_ERR)
		}
		if bucket == "" {
			fmt.Fprintln(os.Stderr, "Missing required param --bucket")
			os.Exit(EXIT_USER_ERR)
		}
		key := cmd.Flags().Lookup("key").Value.String()
		if key == "" {
			key = path.Base(file)
		}

		logger.Infof("Uploading file %s to %s/%s/%s", file, s3Host, bucket, key)
		client := GetS3Client(s3Host)
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
	rootCmd.AddCommand(s3uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// aptrust s3upload source dest
	s3uploadCmd.Flags().StringP("host", "H", "", "S3 host name. E.g. s3.amazonaws.com.")
	s3uploadCmd.Flags().StringP("bucket", "b", "", "Bucket to upload from")
	s3uploadCmd.Flags().StringP("key", "k", "", "Key (name of object) to download")

}
