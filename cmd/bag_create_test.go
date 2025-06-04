package cmd_test

import (
	"fmt"
	"testing"

	"github.com/APTrust/apt-cmd/cmd"
	"github.com/APTrust/dart-runner/bagit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureDefaultTags(t *testing.T) {
	tags := make([]*bagit.TagDefinition, 0)
	tags = cmd.EnsureDefaultTags(tags)
	version := cmd.FindTag(tags, "bagit.txt", "BagIt-Version")
	encoding := cmd.FindTag(tags, "bagit.txt", "Tag-File-Character-Encoding")
	require.NotNil(t, version)
	require.NotNil(t, encoding)
	assert.Equal(t, "1.0", version.GetValue())
	assert.Equal(t, "UTF-8", encoding.GetValue())

	// Make sure that if values are present but empty,
	// we set them correctly without appending additional
	// tags to the list.
	version.UserValue = ""
	encoding.UserValue = ""
	tags = cmd.EnsureDefaultTags(tags)
	version = cmd.FindTag(tags, "bagit.txt", "BagIt-Version")
	encoding = cmd.FindTag(tags, "bagit.txt", "Tag-File-Character-Encoding")
	require.NotNil(t, version)
	require.NotNil(t, encoding)
	require.Equal(t, 2, len(tags))

	// Make sure we don't overwrite values if they're already specified
	version.UserValue = "0.97"
	encoding.UserValue = "ascii"
	tags = cmd.EnsureDefaultTags(tags)
	version = cmd.FindTag(tags, "bagit.txt", "BagIt-Version")
	encoding = cmd.FindTag(tags, "bagit.txt", "Tag-File-Character-Encoding")
	require.NotNil(t, version)
	require.NotNil(t, encoding)
	assert.Equal(t, "0.97", version.GetValue())
	assert.Equal(t, "ascii", encoding.GetValue())

}

func TestValidateTags(t *testing.T) {
	profile, err := cmd.LoadProfile("aptrust")
	require.Nil(t, err)

	expected := []string{
		"Required tag bagit.txt/BagIt-Version is missing.",
		"Required tag bagit.txt/Tag-File-Character-Encoding is missing.",
		"Required tag bag-info.txt/Source-Organization is missing.",
		"Required tag aptrust-info.txt/Title is missing.",
		"Required tag aptrust-info.txt/Access is missing.",
		"Required tag aptrust-info.txt/Storage-Option is missing.",
	}

	tags := make([]*bagit.TagDefinition, 0)
	errors := cmd.ValidateTags(profile, tags)
	assert.Equal(t, len(expected), len(errors))
	assert.Equal(t, expected, errors)

	expected = []string{
		"Required tag bag-info.txt/Source-Organization is missing.",
		"Required tag aptrust-info.txt/Title is missing.",
		"Required tag aptrust-info.txt/Access is missing.",
		"Required tag aptrust-info.txt/Storage-Option is missing.",
	}
	tags = cmd.EnsureDefaultTags(tags)
	errors = cmd.ValidateTags(profile, tags)
	assert.Equal(t, len(expected), len(errors))
	assert.Equal(t, expected, errors)

	tags = append(tags, &bagit.TagDefinition{TagFile: "bag-info.txt", TagName: "Source-Organization", UserValue: "APTrust"})
	tags = append(tags, &bagit.TagDefinition{TagFile: "aptrust-info.txt", TagName: "Title", UserValue: "Bag Title"})
	tags = append(tags, &bagit.TagDefinition{TagFile: "aptrust-info.txt", TagName: "Access", UserValue: "Consortia"})
	tags = append(tags, &bagit.TagDefinition{TagFile: "aptrust-info.txt", TagName: "Storage-Option", UserValue: "Standard"})

	tags = cmd.EnsureDefaultTags(tags)
	errors = cmd.ValidateTags(profile, tags)
	assert.Equal(t, 0, len(errors))

	// Test an illegal value and a value that can't be empty
	tags[3].UserValue = ""        // Empty Title
	tags[4].UserValue = "invalid" // invalid value for Access tag

	expected = []string{
		"Tag aptrust-info.txt/Title is present but value cannot be empty. Please assign a value.",
		"Tag aptrust-info.txt/Access assigned illegal value 'invalid'. Valid values are: Consortia,Institution,Restricted.",
	}

	tags = cmd.EnsureDefaultTags(tags)
	errors = cmd.ValidateTags(profile, tags)
	assert.Equal(t, len(expected), len(errors))
	assert.Equal(t, expected, errors)
}

func TestFindTag(t *testing.T) {
	tags := make([]*bagit.TagDefinition, 0)
	version := cmd.FindTag(tags, "bagit.txt", "BagIt-Version")
	assert.Nil(t, version)
	tags = append(tags, &bagit.TagDefinition{TagFile: "bagit.txt", TagName: "BagIt-Version"})
	version = cmd.FindTag(tags, "bagit.txt", "BagIt-Version")
	assert.NotNil(t, version)
	assert.Equal(t, "BagIt-Version", version.TagName)
}

func TestValidateManifestAlgorithms(t *testing.T) {
	profile, err := cmd.LoadProfile("aptrust")
	require.Nil(t, err)
	require.NotNil(t, profile)

	expected := []string{
		"Profile APTrust requires manifest algorithm md5",
	}

	algs := make([]string, 0)
	errors := cmd.ValidateManifestAlgorithms(profile, algs)
	assert.Equal(t, len(expected), len(errors))
	assert.Equal(t, expected, errors)

	// This should produce no errors, as both these algorithms
	// are allowed by the profile.
	algs = []string{
		"md5",
		"sha256",
	}
	errors = cmd.ValidateManifestAlgorithms(profile, algs)
	fmt.Println(errors)
	assert.Equal(t, 0, len(errors))

	algs = []string{
		"md4",
		"sha1",
		"md5",
		"sha256",
		"sha512",
	}
	expected = []string{
		"Manifest algorithm 'md4' is not allowed in profile APTrust.",
	}
	errors = cmd.ValidateManifestAlgorithms(profile, algs)
	assert.Equal(t, len(expected), len(errors))
	assert.Equal(t, expected, errors)

}
