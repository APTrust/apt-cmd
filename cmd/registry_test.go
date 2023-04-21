//go:build integration && !windows

package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/APTrust/apt-cmd/cmd"
	"github.com/APTrust/preservation-services/models/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: To run integration tests, run `scripts/test.rb integration`.
// That starts up a local instance of registry and loads the same
// test fixtures that we use in registry and preservation-services
// tests.

type FilesResponse struct {
	Count    int                     `json:"count"`
	Next     *string                 `json:"next"`
	Previous *string                 `json:"previous"`
	Results  []*registry.GenericFile `json:"results"`
}

type ObjectsResponse struct {
	Count    int                            `json:"count"`
	Next     *string                        `json:"next"`
	Previous *string                        `json:"previous"`
	Results  []*registry.IntellectualObject `json:"results"`
}

type WorkItemsResponse struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []*registry.WorkItem `json:"results"`
}

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
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "list", "files", "intellectual_object_id=3", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp := &FilesResponse{}
	err := json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	assert.Equal(t, 5, len(resp.Results))

	exitCode, stdout, stderr = execCmd(t, "go", "run", "../main.go", "registry", "list", "files", "intellectual_object_id=3", "sort=identifier__desc", "state=A", "per_page=2", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	//fmt.Println(stdout)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp = &FilesResponse{}
	err = json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	files := resp.Results
	assert.Equal(t, 2, len(files))
	assert.Equal(t, "institution1.edu/glass/shard5-pending-restoration", files[0].Identifier)
	assert.Equal(t, "institution1.edu/glass/shard3", files[1].Identifier)
}

func TestRegistryObjectGet(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "get", "object", "id=1", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	obj := &registry.IntellectualObject{}
	err := json.Unmarshal([]byte(stdout), obj)
	require.Nil(t, err)

	assert.Equal(t, int64(1), obj.ID)
	assert.Equal(t, "institution1.edu/photos", obj.Identifier)
}

func TestRegistryObjectList(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "list", "objects", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp := &ObjectsResponse{}
	err := json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	assert.Equal(t, 6, len(resp.Results))

	exitCode, stdout, stderr = execCmd(t, "go", "run", "../main.go", "registry", "list", "objects", "storage_option=Standard", "sort=identifier__desc", "state=A", "per_page=2", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	//fmt.Println(stdout)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp = &ObjectsResponse{}
	err = json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	objects := resp.Results
	assert.Equal(t, 2, len(objects))
	assert.Equal(t, "institution1.edu/photos", objects[0].Identifier)
	assert.Equal(t, "institution1.edu/pdfs", objects[1].Identifier)
}

func TestRegistryWorkItemGet(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "get", "workitem", "id=25", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	item := &registry.WorkItem{}
	err := json.Unmarshal([]byte(stdout), item)
	require.Nil(t, err)

	assert.Equal(t, int64(25), item.ID)
	assert.Equal(t, "institution1.edu/photos", item.ObjectIdentifier)
	assert.Equal(t, "Cleanup", item.Stage)
	assert.Equal(t, "Success", item.Status)
}

func TestRegistryWorkItemList(t *testing.T) {
	exitCode, stdout, stderr := execCmd(t, "go", "run", "../main.go", "registry", "list", "workitems", "action=Ingest", "stage=Receive", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp := &WorkItemsResponse{}
	err := json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	assert.Equal(t, 10, len(resp.Results))

	// Limit to 2. Will order by ID if we don't specify a sort param.
	exitCode, stdout, stderr = execCmd(t, "go", "run", "../main.go", "registry", "list", "workitems", "action=Ingest", "stage=Cleanup", "per_page=2", "--config=../testconfig.env")
	assert.Equal(t, cmd.EXIT_OK, exitCode)
	//fmt.Println(stdout)
	require.True(t, len(stdout) > 100)
	assert.Equal(t, "", stderr)

	resp = &WorkItemsResponse{}
	err = json.Unmarshal([]byte(stdout), resp)
	require.Nil(t, err)
	items := resp.Results
	assert.Equal(t, 2, len(items))
	assert.Equal(t, int64(22), items[0].ID)
	assert.Equal(t, int64(23), items[1].ID)
}
