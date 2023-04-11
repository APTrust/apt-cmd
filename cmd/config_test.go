package cmd_test

import (
	"testing"

	"github.com/APTrust/partner-tools/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig(populated bool) *cmd.Config {
	if !populated {
		return &cmd.Config{}
	}
	return &cmd.Config{
		RegistryURL:        "https://demo.aptrust.org",
		RegistryAPIVersion: "v3",
		RegistryEmail:      "user@example.com",
		RegistryAPIKey:     "top-seekrit!",
		AWSKey:             "AWS-KEY-1",
		AWSSecret:          "AWS-SECRET-1",
		ConfigSource:       "getTestConfig",
	}
}

func TestValidateRegistryConfig(t *testing.T) {
	emptyConfig := getTestConfig(false)
	err := emptyConfig.ValidateRegistryConfig()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "URL is missing")
	assert.Contains(t, err.Error(), "API version is missing")
	assert.Contains(t, err.Error(), "email is missing")
	assert.Contains(t, err.Error(), "API key is missing")

	fullConfig := getTestConfig(true)
	err = fullConfig.ValidateRegistryConfig()
	assert.Nil(t, err)
}

func TesValidateAWSCredentials(t *testing.T) {
	emptyConfig := getTestConfig(false)
	err := emptyConfig.ValidateAWSCredentials()
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key is missing")
	assert.Contains(t, err.Error(), "Secret is missing")

	fullConfig := getTestConfig(true)
	err = fullConfig.ValidateAWSCredentials()
	assert.Nil(t, err)
}

func TestConfigString(t *testing.T) {
	expectedEmpty := `Configuration:
	RegistryURL:             
	RegistryAPIVersion:      
	RegistryEmail:           
	RegistryAPIKey:          MISSING!
	AWSKey:                  MISSING!
	AWSSecret:               MISSING!
	ConfigSource:            `
	emptyConfig := getTestConfig(false)
	assert.Equal(t, expectedEmpty, emptyConfig.String())

	expecteFull := `Configuration:
	RegistryURL:             https://demo.aptrust.org
	RegistryAPIVersion:      v3
	RegistryEmail:           user@example.com
	RegistryAPIKey:          [redacted]
	AWSKey:                  [redacted]
	AWSSecret:               [redacted]
	ConfigSource:            getTestConfig`
	fullConfig := getTestConfig(true)
	assert.Equal(t, expecteFull, fullConfig.String())
}
