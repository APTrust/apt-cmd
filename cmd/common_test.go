package cmd_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/APTrust/apt-cmd/cmd"
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
	emptyConfig := getTestConfig(false)
	client, err := cmd.NewRegistryClient(emptyConfig)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	config := getTestConfig(true)
	client, err = cmd.NewRegistryClient(config)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, config.RegistryAPIKey, client.APIKey)
	assert.Equal(t, config.RegistryAPIVersion, client.APIVersion)
	assert.Equal(t, config.RegistryEmail, client.APIUser)
	assert.Equal(t, config.RegistryURL, client.HostURL)
}

func TestInitRegistryRequest(t *testing.T) {
	config := getTestConfig(true)
	client, urlValues := cmd.InitRegistryRequest(config, testArgs)
	assert.NotNil(t, client)
	assert.Equal(t, 3, len(urlValues))
	assert.Equal(t, "value1", urlValues.Get("pair1"))
	assert.Equal(t, "value2", urlValues.Get("pair2"))
	assert.Equal(t, "value3", urlValues.Get("pair3"))
}

func TestEnsureDefaultListParams(t *testing.T) {
	params := url.Values{}
	assert.Empty(t, params.Get("sort"))
	assert.Empty(t, params.Get("per_page"))

	cmd.EnsureDefaultListParams(params)
	assert.Equal(t, "id", params.Get("sort"))
	assert.Equal(t, "25", params.Get("per_page"))

	// This function should not overwrite param
	// values if they're already specified.
	params.Set("sort", "name__desc")
	params.Set("per_page", "100")
}

func TestPrettyPrintJSON(t *testing.T) {
	uglyAssJson := []byte(`{"id":1,"name":"Oscar","home":"Garbage Can"}`)
	expected := `{
  "id": 1,
  "name": "Oscar",
  "home": "Garbage Can"
}
`
	tempFile, err := os.CreateTemp(os.TempDir(), "ptool-test")
	require.Nil(t, err)

	actualStdOut := os.Stdout
	os.Stdout = tempFile
	defer func() {
		os.Stdout = actualStdOut
	}()

	cmd.PrettyPrintJSON(uglyAssJson)
	tempFile.Close()

	prettyJson, err := ioutil.ReadFile(tempFile.Name())
	require.Nil(t, err)

	assert.Equal(t, expected, string(prettyJson))
}

func TestNewS3Client(t *testing.T) {
	config := getTestConfig(true)
	client := cmd.NewS3Client(config, "s3.amazonaws.com")
	assert.NotNil(t, client)
	assert.Equal(t, "s3.amazonaws.com", client.EndpointURL().Host)
}

func TestLooksLikePreservationBucket(t *testing.T) {
	yes := []string{
		"aptrust.preservation.oregon",
		"aptrust.preservation.storage",
		"aptrust.test.preservation.va",
		"aptrust-staging-wasabi-va",
		"aptrust-wasabi-or",
	}
	no := []string{
		"aptrust.receiving.virginia.edu",
		"aptrust.test.restore.vt.edu",
		"apt-wasabi-staging-va",
		"apt-wasabi-test-or",
	}
	for _, bucket := range yes {
		assert.True(t, cmd.LooksLikePreservationBucket(bucket))
	}
	for _, bucket := range no {
		assert.False(t, cmd.LooksLikePreservationBucket(bucket))
	}
}

func TestGetParam(t *testing.T) {
	// Currently, there's no way to test this, except by
	// using exec and hacking something ugly :(
}

func TestLoadProfile(t *testing.T) {
	profile, err := cmd.LoadProfile("aptrust")
	require.Nil(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "https://raw.githubusercontent.com/APTrust/preservation-services/master/profiles/aptrust-v2.3.json", profile.BagItProfileInfo.BagItProfileIdentifier)

	profile, err = cmd.LoadProfile("aptrust-2.3")
	require.Nil(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "https://raw.githubusercontent.com/APTrust/preservation-services/master/profiles/aptrust-v2.3.json", profile.BagItProfileInfo.BagItProfileIdentifier)

	profile, err = cmd.LoadProfile("aptrust-2.2")
	require.Nil(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "https://raw.githubusercontent.com/APTrust/preservation-services/master/profiles/aptrust-v2.2.json", profile.BagItProfileInfo.BagItProfileIdentifier)

	profile, err = cmd.LoadProfile("btr")
	require.Nil(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "https://raw.githubusercontent.com/dpscollaborative/btr_bagit_profile/master/btr-bagit-profile.json", profile.BagItProfileInfo.BagItProfileIdentifier)

	profile, err = cmd.LoadProfile("empty")
	require.Nil(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, "https://raw.githubusercontent.com/APTrust/dart/tree/master/profiles/empty_profile.json", profile.BagItProfileInfo.BagItProfileIdentifier)
}
