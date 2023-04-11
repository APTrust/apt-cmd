package cmd_test

import (
	"testing"

	"github.com/APTrust/partner-tools/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// For testing exit codes, see
// https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go
// https://github.com/bouk/monkey

var testArgs = []string{
	"arg1",
	"--flag1",
	"pair1=value1",
	"arg2",
	"pair2=value2",
	"--flag3",
	"pair3=value3",
}

func TestParseArgPairs(t *testing.T) {
	pairs := cmd.ParseArgPairs(testArgs)
	require.Equal(t, 3, len(pairs))
	assert.Equal(t, "pair1", pairs[0].Name)
	assert.Equal(t, "value1", pairs[0].Value)
	assert.Equal(t, "pair2", pairs[1].Name)
	assert.Equal(t, "value2", pairs[1].Value)
	assert.Equal(t, "pair3", pairs[2].Name)
	assert.Equal(t, "value3", pairs[2].Value)
}

func TestGetUrlValues(t *testing.T) {
	values := cmd.GetUrlValues(testArgs)
	require.Equal(t, 3, len(values))
	assert.Equal(t, "value1", values.Get("pair1"))
	assert.Equal(t, "value2", values.Get("pair2"))
	assert.Equal(t, "value3", values.Get("pair3"))
}

func TestGetTagValues(t *testing.T) {
	tagArgs := append(testArgs, "aptrust-info.txt/Title=Bag Title")
	tagArgs = append(tagArgs, "random-info.txt/Quarterback=Jim Plunkett")
	tagArgs = append(tagArgs, "bag-info.txt/state-name=Virginia")
	tags := cmd.GetTagValues(tagArgs)
	require.Equal(t, 6, len(tags))

	// Note that tag names are converted to title case
	assert.Equal(t, "bag-info.txt", tags[0].TagFile)
	assert.Equal(t, "Pair1", tags[0].TagName)
	assert.Equal(t, "value1", tags[0].UserValue)
	assert.Equal(t, "bag-info.txt", tags[1].TagFile)
	assert.Equal(t, "Pair2", tags[1].TagName)
	assert.Equal(t, "value2", tags[1].UserValue)
	assert.Equal(t, "bag-info.txt", tags[2].TagFile)
	assert.Equal(t, "Pair3", tags[2].TagName)
	assert.Equal(t, "value3", tags[2].UserValue)

	assert.Equal(t, "aptrust-info.txt", tags[3].TagFile)
	assert.Equal(t, "Title", tags[3].TagName)
	assert.Equal(t, "Bag Title", tags[3].UserValue)

	assert.Equal(t, "random-info.txt", tags[4].TagFile)
	assert.Equal(t, "Quarterback", tags[4].TagName)
	assert.Equal(t, "Jim Plunkett", tags[4].UserValue)

	assert.Equal(t, "bag-info.txt", tags[5].TagFile)
	assert.Equal(t, "State-Name", tags[5].TagName)
	assert.Equal(t, "Virginia", tags[5].UserValue)
}

func TestNewRegistryClient(t *testing.T) {

}

func TestInitRegistryRequest(t *testing.T) {

}

func TestEnsureDefaultListParams(t *testing.T) {

}

func TestPrettyPrintJSON(t *testing.T) {

}

func TestGetS3Client(t *testing.T) {

}

func TestDisallowPreservationBucket(t *testing.T) {

}

func TestGetParam(t *testing.T) {

}
