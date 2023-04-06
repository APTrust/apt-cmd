package cmd

import "fmt"

type Config struct {
	RegistryURL        string
	RegistryAPIVersion string
	RegistryEmail      string
	RegistryAPIKey     string
	AWSKey             string
	AWSSecret          string
	ReceivingBucket    string
	RestorationBucket  string
	DownloadDir        string
	ConfigSource       string
}

func (config *Config) ValidateRegistryConfig() error {
	errMsg := ""
	if config.RegistryURL == "" {
		errMsg = "Registry URL is missing from config. "
	}
	if config.RegistryAPIVersion == "" {
		errMsg += "Registry API version is missing from config. "
	}
	if config.RegistryEmail == "" {
		errMsg += "Registry email is missing from config. "
	}
	if config.RegistryAPIKey == "" {
		errMsg += "Registry API key is missing from config. "
	}
	if errMsg != "" {
		errMsg += fmt.Sprintf("(Config source: %s)", config.ConfigSource)
		return fmt.Errorf(errMsg)
	}
	return nil
}

func (config *Config) ValidateAWSCredentials() error {
	errMsg := ""
	if config.AWSKey == "" {
		errMsg = "AWS Key is missing from config. "
	}
	if config.AWSSecret == "" {
		errMsg += "AWS Secret is missing from config. "
	}
	if errMsg != "" {
		errMsg += fmt.Sprintf("(Config source: %s)", config.ConfigSource)
		return fmt.Errorf(errMsg)
	}
	return nil
}