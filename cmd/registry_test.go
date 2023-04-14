//go:build integration
// +build integration

package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/APTrust/partner-tools/cmd"
	"github.com/APTrust/preservation-services/models/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var intTestConfig = &cmd.Config{
	RegistryURL:        "http://localhost:8080",
	RegistryEmail:      "user@inst1.edu",
	RegistryAPIVersion: "v3",
	RegistryAPIKey:     "password",
	AWSKey:             "minioadmin",
	AWSSecret:          "minioadmin",
}

func TestRegistryFileGet(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "get", "file", "id=1", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	gf := &registry.GenericFile{}
	err := json.Unmarshal([]byte(stdout), gf)
	require.Nil(t, err)

	assert.Equal(t, int64(1), gf.ID)
	assert.Equal(t, "institution1.edu/photos/picture1", gf.Identifier)
	assert.Equal(t, 2, len(gf.StorageRecords))
	assert.Equal(t, 2, len(gf.Checksums))
	assert.Equal(t, 4, len(gf.PremisEvents))
}

func TestRegistryFileList(t *testing.T) {

}
