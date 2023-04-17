package cmd_test

import (
	"fmt"
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

}

func testS3Delete(t *testing.T) {

}
