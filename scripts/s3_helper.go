//go:build ignore

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	Bucket    = "aptrust.public.download"
	KeyPrefix = "apt-cmd"
	Endpoint  = "s3.amazonaws.com"
)

func main() {
	makeFolders := flag.Bool("make-folders", false, "Create S3 folders for a release version")
	upload := flag.Bool("upload", false, "Upload a binary file to S3")
	list := flag.Bool("list", false, "List S3 objects under a prefix")
	getLinks := flag.String("get-links", "", "Get public download links for a specific version (e.g. v3.0.4)")
	version := flag.String("version", "", "apt-cmd version (e.g. v3.0.4)")
	arch := flag.String("arch", "", "OS/arch path component (e.g. linux/amd64, mac/arm64, windows/amd64)")
	prefix := flag.String("prefix", "", "Prefix for the -list command")

	flag.Parse()

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set in the environment")
	}

	s3Client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Error creating S3 client: %v", err)
	}

	success := true
	switch {
	case *makeFolders:
		if *version == "" {
			log.Fatal("-version is required for -make-folders")
		}
		success = doMakeFolders(s3Client, *version)

	case *upload:
		if *version == "" || *arch == "" {
			log.Fatal("-version and -arch are required for -upload")
		}
		args := flag.Args()
		if len(args) < 1 {
			log.Fatal("please provide the path to the file to upload as a positional argument")
		}
		success = doUpload(s3Client, *version, *arch, args[0])

	case *list:
		if *prefix == "" {
			log.Fatal("-prefix is required for -list")
		}
		success = doList(s3Client, *prefix)

	case *getLinks != "":
		success = doGetLinks(s3Client, *getLinks)

	default:
		fmt.Fprintln(os.Stderr, "Please specify an action: -make-folders, -upload, -list, or -get-links")
		flag.Usage()
		success = false
	}

	if !success {
		os.Exit(1)
	}
}

func doMakeFolders(client *minio.Client, version string) bool {
	ctx := context.Background()
	arches := []string{
		"linux/amd64",
		"linux/arm64",
		"mac/amd64",
		"mac/arm64",
		"windows/amd64",
		"windows/arm64",
	}
	success := true

	for _, a := range arches {
		folderKey := fmt.Sprintf("%s/%s/%s/", KeyPrefix, version, a)

		_, err := client.StatObject(ctx, Bucket, folderKey, minio.StatObjectOptions{})
		if err == nil {
			fmt.Printf("Folder already exists: s3://%s/%s\n", Bucket, folderKey)
			continue
		}

		_, err = client.PutObject(ctx, Bucket, folderKey, bytes.NewReader([]byte{}), 0, minio.PutObjectOptions{
			UserMetadata: map[string]string{"x-amz-acl": "public-read"},
		})
		if err != nil {
			fmt.Printf("Failed to create folder s3://%s/%s: %v\n", Bucket, folderKey, err)
			success = false
		} else {
			fmt.Printf("Created folder: s3://%s/%s\n", Bucket, folderKey)
		}
	}
	return success
}

func doUpload(client *minio.Client, version, arch, filePath string) bool {
	ctx := context.Background()
	fileName := filepath.Base(filePath)
	// arch already contains the os/arch path component, e.g. "linux/amd64"
	objectKey := fmt.Sprintf("%s/%s/%s/%s", KeyPrefix, version, arch, fileName)

	// Calculate SHA256 checksum before uploading
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", filePath, err)
		return false
	}
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		fmt.Printf("Failed to calculate SHA256 for %s: %v\n", filePath, err)
		f.Close()
		return false
	}
	f.Close()
	checksum := hex.EncodeToString(h.Sum(nil))

	_, err = client.FPutObject(ctx, Bucket, objectKey, filePath, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		UserMetadata: map[string]string{
			"x-amz-acl":         "public-read",
			"x-amz-meta-sha256": checksum,
		},
	})
	if err != nil {
		fmt.Printf("Failed to upload %s to s3://%s/%s: %v\n", filePath, Bucket, objectKey, err)
		return false
	}

	fmt.Printf("Uploaded %s -> s3://%s/%s\n", filePath, Bucket, objectKey)
	fmt.Printf("SHA256: %s\n", checksum)
	return true
}

func doList(client *minio.Client, prefix string) bool {
	ctx := context.Background()

	fmt.Printf("%-24s %-15s %s\n", "Last Updated", "Size (bytes)", "Object Key")
	fmt.Printf("%-24s %-15s %s\n", strings.Repeat("-", 24), strings.Repeat("-", 15), strings.Repeat("-", 40))

	for object := range client.ListObjects(ctx, Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if object.Err != nil {
			fmt.Printf("List error: %v\n", object.Err)
			return false
		}
		dateStr := object.LastModified.Format("2006-01-02 15:04:05 MST")
		fmt.Printf("%-24s %-15d %s\n", dateStr, object.Size, object.Key)
	}
	return true
}

func doGetLinks(client *minio.Client, version string) bool {
	ctx := context.Background()
	prefix := fmt.Sprintf("%s/%s", KeyPrefix, version)

	fmt.Printf("%-24s %-15s %-64s %s\n", "Last Updated", "Size (bytes)", "SHA256 Checksum", "Download Link")
	fmt.Printf("%-24s %-15s %-64s %s\n",
		strings.Repeat("-", 24), strings.Repeat("-", 15), strings.Repeat("-", 64), strings.Repeat("-", 15))

	for object := range client.ListObjects(ctx, Bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if object.Err != nil {
			fmt.Printf("List error: %v\n", object.Err)
			return false
		}
		if object.Size == 0 {
			continue // skip folder markers
		}

		link := fmt.Sprintf("https://%s/%s/%s", Endpoint, Bucket, object.Key)
		checksum := "N/A"

		objInfo, err := client.StatObject(ctx, Bucket, object.Key, minio.StatObjectOptions{})
		if err == nil {
			for _, key := range []string{"Sha256", "x-amz-meta-sha256", "X-Amz-Meta-Sha256"} {
				if val := objInfo.UserMetadata[key]; val != "" {
					checksum = val
					break
				}
			}
			if checksum == "N/A" {
				if val := objInfo.Metadata.Get("X-Amz-Meta-Sha256"); val != "" {
					checksum = val
				}
			}
		}

		dateStr := object.LastModified.Format("2006-01-02 15:04:05 MST")
		fmt.Printf("%-24s %-15d %-64s %s\n", dateStr, object.Size, checksum, link)
	}
	return true
}
