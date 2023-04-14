package cmd_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/APTrust/partner-tools/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBagCreate(t *testing.T) {
	tmpFile := path.Join("..", "partnertools-testbag.tar")
	defer os.Remove(tmpFile)

	args := []string{
		"run",
		"../main.go",
		"bag",
		"create",
		`--profile=aptrust`,
		`--manifest-algs=md5,sha256`,
		fmt.Sprintf(`--output-file=%s`, tmpFile),
		`--bag-dir=profiles`,
		`--tags=aptrust-info.txt/Title=Bag of Profiles`,
		`--tags=aptrust-info.txt/Access=Institution`,
		`--tags=aptrust-info.txt/Storage-Option=Standard`,
		`--tags=bag-info.txt/Source-Organization=Faber College`,
		`--tags=Custom-Tag=Single quoted because it {contains} $weird &characters`,
	}

	exitCode, stdout, stderr := execCmd(t, "go", args...)
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "", stderr)
	assert.Contains(t, stdout, `"result": "OK"`)
	assert.Contains(t, stdout, "partnertools-testbag.tar") // Tells us where the bag is

	// Make sure that bag is valid
	if exitCode == 0 {
		exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "bag", "validate", "--profile=aptrust", tmpFile)
		assert.Equal(t, 0, exitCode)
		assert.Equal(t, "Bag is valid according to aptrust profile.\n", stdout)
		assert.Equal(t, "", stderr)
	}
}

func TestBagValidate_GoodBags(t *testing.T) {
	goodAPTrustBags := []string{
		"example.edu.sample_good.tar",
		"example.edu.tagsample_good.tar",
	}
	goodBTRBags := []string{
		"test.edu.btr-glacier-deep-oh.tar",
		"test.edu.btr_good_sha256.tar",
		"test.edu.btr_good_sha512.tar",
		"test.edu.btr-wasabi-or.tar",
	}

	for _, bag := range goodAPTrustBags {
		testValidBag(t, "aptrust", bag)
	}

	for _, bag := range goodBTRBags {
		testValidBag(t, "btr", bag)
	}
}

func TestBagValidate_BadBags(t *testing.T) {
	testInvalidBag(t, "btr", "test.edu.btr_bad_checksums.tar", "does not match digest")
	testInvalidBag(t, "btr", "test.edu.btr_bad_extraneous_file.tar", "Payload-Oxum does not match payload")
	testInvalidBag(t, "btr", "test.edu.btr_bad_missing_payload_file.tar", "Payload-Oxum does not match payload")
	testInvalidBag(t, "btr", "test.edu.btr_bad_missing_required_tags.tar", "Required tag is missing")

	testInvalidBag(t, "aptrust", "example.edu.sample_bad_oxum.tar", "Payload-Oxum does not match payload")
	testInvalidBag(t, "aptrust", "example.edu.sample_missing_data_file.tar", "file is missing from bag")
	testInvalidBag(t, "aptrust", "example.edu.sample_no_bag_info.tar", "Required tag is missing")
	testInvalidBag(t, "aptrust", "example.edu.sample_no_data_dir.tar", "file is missing from bag", "Required tag is missing")
	testInvalidBag(t, "aptrust", "example.edu.sample_no_md5_manifest.tar", "Required manifest is missing", "Required tag is missing")
	testInvalidBag(t, "aptrust", "example.edu.tagsample_bad.tar", "Tag has illegal value", "Required tag is present but has no value", "does not match digest", "file is missing from bag")
}

func testValidBag(t *testing.T, profileName, tarFileName string) {
	profileFlag := fmt.Sprintf("--profile=%s", profileName)
	pathToBag := path.Join("..", "testbags", profileName, tarFileName)
	expectedStdout := fmt.Sprintf("Bag is valid according to %s profile.\n", profileName)
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "bag", "validate", profileFlag, pathToBag)
	assert.Equal(t, 0, exitCode, tarFileName)
	assert.Equal(t, expectedStdout, stdout, tarFileName)
	assert.Equal(t, "", stderr, tarFileName)
}

// Note: When testing invalid bags, validation messages go to stdout, not stderr,
// because the validator ran to completion. Also note that even though the exit
// status for an invalid bag is 2, go's cmd.Wait returns 1. Not sure why, but we
// do see the expected 2 in the output from stderr.
func testInvalidBag(t *testing.T, profileName, tarFileName string, errorMsgSubstring ...string) {
	profileFlag := fmt.Sprintf("--profile=%s", profileName)
	pathToBag := path.Join("..", "testbags", profileName, tarFileName)
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "bag", "validate", profileFlag, pathToBag)
	assert.NotEqual(t, cmd.EXIT_OK, exitCode, tarFileName)
	for _, msg := range errorMsgSubstring {
		assert.Contains(t, stdout, msg, tarFileName)
	}
	assert.Equal(t, "exit status 2\n", stderr, tarFileName)
}

// execCmd runs a command and return the exit code, stdout and stderr output.
// We run tests through exec rather than in-process because all of these
// commands call os.Exit() deliberately if something is wrong. We don't want
// our test process to exit.
func execCmd(t *testing.T, commandName string, args ...string) (int, string, string) {
	cmd := exec.Command(commandName, args...)
	stdout, err := cmd.StdoutPipe()
	require.Nil(t, err)
	stderr, err := cmd.StderrPipe()
	require.Nil(t, err)

	err = cmd.Start()
	require.Nil(t, err)

	stdoutData, err := io.ReadAll(stdout)
	require.Nil(t, err)
	stderrData, err := io.ReadAll(stderr)
	require.Nil(t, err)

	// Don't check the error here, because in some cases
	// the caller is expecting a specific exit error code.
	// We'll return that and let the caller check it.
	cmd.Wait()

	return cmd.ProcessState.ExitCode(), string(stdoutData), string(stderrData)
}
