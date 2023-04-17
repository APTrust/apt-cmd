//go:build integration
// +build integration

package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/APTrust/partner-tools/cmd"
	"github.com/stretchr/testify/assert"
)

// localhost:9899

var s3TestFiles = []string{
	"bag_create.go",
	"bag.go",
	"common.go",
	"config.go",
}

func TestS3(t *testing.T) {
	testS3Upload(t)
	testS3List(t)
	testS3Download(t)
	testS3Delete(t)
}

func testS3Upload(t *testing.T) {
	for _, file := range s3TestFiles {
		exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "s3", "upload", "--host=127.0.0.1:9899", "--bucket=test-bucket-1", "--config=../testconfig.env", file)
		assert.Equal(t, cmd.EXIT_OK, exitCode)
		assert.NotEmpty(t, stdout)
		assert.Empty(t, stderr)

		// Make sure stdout shows JSON for this file upload.
		key := fmt.Sprintf(`"Key": "%s"`, file)
		assert.Contains(t, stdout, key)
	}
}

func testS3List(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "s3", "list", "--host=127.0.0.1:9899", "--bucket=test-bucket-1", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	assert.NotEmpty(t, stdout)
	assert.Empty(t, stderr)

	for _, file := range s3TestFiles {
		assert.Contains(t, stdout, file)
	}
}

func testS3Download(t *testing.T) {
	defer os.Remove("download-test.txt")
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "s3", "download", "--host=127.0.0.1:9899", "--bucket=test-bucket-1", "--key=bag.go", "--save-as=download-test.txt", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "download-test.txt")
}

func testS3Delete(t *testing.T) {
	for _, file := range s3TestFiles {
		exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "s3", "delete", "--host=127.0.0.1:9899", "--bucket=test-bucket-1", "--config=../testconfig.env", "--key="+file)
		assert.Equal(t, cmd.EXIT_OK, exitCode)
		assert.NotEmpty(t, stdout)
		assert.Empty(t, stderr)
		assert.Contains(t, stdout, fmt.Sprintf("test-bucket-1/%s", file))
	}
}
