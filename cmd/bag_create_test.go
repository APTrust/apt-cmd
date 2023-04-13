package cmd_test

import (
	"testing"

	"github.com/APTrust/dart-runner/bagit"
	"github.com/APTrust/partner-tools/cmd"
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
